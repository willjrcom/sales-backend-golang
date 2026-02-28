package orderusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type OrderProcessService struct {
	r   model.OrderProcessRepository
	rpr model.ProcessRuleRepository
	sq  *orderqueueusecases.Service
	sgi *GroupItemService
	rgi model.GroupItemRepository
	ro  model.OrderRepository
	se  *employeeusecases.Service
	so  *OrderService
}

func NewOrderProcessService(c model.OrderProcessRepository) *OrderProcessService {
	return &OrderProcessService{r: c}
}

func (s *OrderProcessService) AddDependencies(sq *orderqueueusecases.Service, rpr model.ProcessRuleRepository, sgi *GroupItemService, ro model.OrderRepository, se *employeeusecases.Service, rgi model.GroupItemRepository, so *OrderService) {
	s.rgi = rgi
	s.rpr = rpr
	s.sq = sq
	s.sgi = sgi
	s.ro = ro
	s.se = se
	s.so = so
}

func (s *OrderProcessService) CreateProcess(ctx context.Context, dto *orderprocessdto.OrderProcessCreateDTO) (uuid.UUID, error) {
	process, err := dto.ToDomain()
	if err != nil {
		return uuid.Nil, err
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, process.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, err
	}

	groupItem := groupItemModel.ToDomain()
	process.OrderID = groupItem.OrderID

	productIDs, err := groupItem.GetDistinctProductIDs()
	if err != nil {
		return uuid.Nil, err
	}

	for _, productID := range productIDs {
		process.Products = append(process.Products, productentity.Product{
			Entity: entity.Entity{
				ID: productID,
			},
		})
	}

	processModel := &model.OrderProcess{}
	processModel.FromDomain(process)
	if err := s.r.CreateProcess(ctx, processModel); err != nil {
		return uuid.Nil, err
	}

	// producerKafka := kafka.NewProducer()
	// if err := producerKafka.NewMessage(ctx, "order_process", process); err != nil {
	// 	return uuid.Nil, err
	// }

	return process.ID, nil
}

func (s *OrderProcessService) StartProcess(ctx context.Context, dtoID *entitydto.IDRequest) error {
	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	employee, err := s.se.GetEmployeeByUserID(ctx, entitydto.NewIdRequest(userIDUUID))
	if err != nil {
		return err
	}

	processModel, err := s.r.GetProcessById(ctx, dtoID.ID.String(), false)
	if err != nil {
		return err
	}

	process := processModel.ToDomain()

	processRule, err := s.rpr.GetProcessRuleById(ctx, process.ProcessRuleID.String())
	if err != nil {
		return err
	}

	if processRule.Order == 1 {
		entityDtoID := entitydto.NewIdRequest(process.GroupItemID)
		if err := s.sgi.StartGroupItem(ctx, entityDtoID); err != nil {
			fmt.Printf("error starting group item: %v", err)
		}
	}

	if err := process.StartProcess(employee.ID); err != nil {
		return err
	}

	processModel.FromDomain(process)
	if err := s.r.UpdateProcess(ctx, processModel); err != nil {
		return err
	}

	if err := s.sq.FinishQueue(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *OrderProcessService) FinishProcess(ctx context.Context, dtoID *entitydto.IDRequest) (nextProcessID uuid.UUID, err error) {
	processModel, err := s.r.GetProcessById(ctx, dtoID.ID.String(), false)
	if err != nil {
		return uuid.Nil, err
	}

	process := processModel.ToDomain()
	if err := process.FinishProcess(); err != nil {
		return uuid.Nil, err
	}

	processModel.FromDomain(process)
	if err := s.r.UpdateProcess(ctx, processModel); err != nil {
		return uuid.Nil, err
	}

	isLast, err := s.rpr.IsLastProcessRuleByID(ctx, process.ProcessRuleID)
	if err != nil {
		return uuid.Nil, err
	}

	// Processes finished
	if isLast {
		entityDtoID := &entitydto.IDRequest{ID: process.GroupItemID}
		if err := s.sgi.ReadyGroupItem(ctx, entityDtoID); err != nil {
			fmt.Printf("error ready group item: %v", err)
		}

		groupItemDTO, err := s.sgi.GetGroupByID(ctx, entityDtoID)
		if err != nil {
			return uuid.Nil, err
		}

		orderModel, err := s.ro.GetOrderById(ctx, groupItemDTO.OrderID.String())
		if err != nil {
			return uuid.Nil, err
		}

		order := orderModel.ToDomain()
		orderIsReady := true

		// Search not ready group item
		for _, groupItem := range order.GroupItems {
			// Next group item
			if groupItem.ID == process.GroupItemID {
				continue
			}

			// Is not ready
			if groupItem.Status != orderentity.StatusGroupReady {
				orderIsReady = false
			}
		}

		if orderIsReady {
			// Update order status only if it's not already Ready or Finished
			if order.Status == orderentity.OrderStatusPending {
				if err := s.so.ReadyOrder(ctx, &entitydto.IDRequest{ID: order.ID}); err != nil {
					return uuid.Nil, err
				}
			}
			return uuid.Nil, nil
		}

		company, err := s.so.sc.GetCompany(ctx)
		if err != nil {
			fmt.Printf("error getting company: %v", err)
			return uuid.Nil, nil
		}

		if enablePrintItemsOnFinishProcess, _ := company.Preferences.GetBool(companyentity.EnablePrintItemsOnFinishProcess); enablePrintItemsOnFinishProcess {
			path := rabbitmq.GROUP_ITEM_PATH + process.GroupItemID.String()
			if err := s.so.rabbitmq.SendPrintMessage(rabbitmq.GROUP_ITEM_EX, company.SchemaName, path, groupItemDTO.PrinterName); err != nil {
				fmt.Printf("error sending message: %v", err)
			}
		}

		// finished last process for this group item, no next process
		return uuid.Nil, nil
	}

	processRule, err := s.rpr.GetProcessRuleById(ctx, process.ProcessRuleID.String())
	if err != nil {
		return uuid.Nil, err
	}

	nextProcessRule, err := s.rpr.GetProcessRuleByCategoryIdAndOrder(ctx, processRule.CategoryID.String(), processRule.Order+1)
	if err != nil {
		return uuid.Nil, err
	}

	// Create next process (idempotent: skip if one already exists)
	existing, err := s.r.GetActiveProcessByGroupItemAndProcessRule(ctx, process.GroupItemID.String(), nextProcessRule.ID.String())
	if err != nil {
		return uuid.Nil, err
	}
	if existing != nil {
		// Process already exists for the next step — return its ID without creating a duplicate
		return existing.ID, nil
	}

	// Start next queue only if the process doesn't exist yet
	startQueueInput := &orderqueuedto.QueueCreateDTO{
		GroupItemID: process.GroupItemID,
		JoinedAt:    *process.FinishedAt,
	}

	if _, err := s.sq.StartQueue(ctx, startQueueInput); err != nil {
		return uuid.Nil, err
	}

	createProcessInput := &orderprocessdto.OrderProcessCreateDTO{
		OrderNumber:   process.OrderNumber,
		OrderType:     process.OrderType,
		GroupItemID:   process.GroupItemID,
		ProcessRuleID: nextProcessRule.ID,
	}

	nextProcessID, err = s.CreateProcess(ctx, createProcessInput)
	if err != nil {
		return uuid.Nil, err
	}

	return nextProcessID, nil
}

func (s *OrderProcessService) CancelProcess(ctx context.Context, dtoID *entitydto.IDRequest, orderprocessdto *orderprocessdto.OrderProcessCancelDTO) error {
	reason, err := orderprocessdto.ToDomain()
	if err != nil {
		return err
	}

	processModel, err := s.r.GetProcessById(ctx, dtoID.ID.String(), false)
	if err != nil {
		return err
	}

	process := processModel.ToDomain()
	if err := process.CancelProcess(reason); err != nil {
		return err
	}

	if err := s.sq.FinishQueue(ctx, process); err != nil {
		return err
	}

	processModel.FromDomain(process)
	if err := s.r.UpdateProcess(ctx, processModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderProcessService) GetProcessById(ctx context.Context, dto *entitydto.IDRequest) (*orderprocessdto.OrderProcessDTO, error) {
	if processModel, err := s.r.GetProcessById(ctx, dto.ID.String(), true); err != nil {
		return nil, err
	} else {
		process := processModel.ToDomain()
		processDTO := &orderprocessdto.OrderProcessDTO{}
		processDTO.FromDomain(process)
		return processDTO, nil
	}
}

func (s *OrderProcessService) GetAllProcessesFinishedByShiftID(ctx context.Context, dto *entitydto.IDRequest) ([]orderprocessdto.OrderProcessDTO, error) {
	if processModels, err := s.r.GetAllProcessesFinishedByShiftID(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return s.modelsToDTOs(processModels), nil
	}
}

func (s *OrderProcessService) GetProcessesByProcessRuleID(ctx context.Context, dtoID *entitydto.IDRequest) ([]orderprocessdto.OrderProcessDTO, error) {
	if processModels, err := s.r.GetProcessesByProcessRuleID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.modelsToDTOs(processModels), nil
	}
}

func (s *OrderProcessService) GetProcessesByProductID(ctx context.Context, dtoID *entitydto.IDRequest) ([]orderprocessdto.OrderProcessDTO, error) {
	if processModels, err := s.r.GetProcessesByProductID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.modelsToDTOs(processModels), nil
	}
}

func (s *OrderProcessService) GetProcessesByGroupItemID(ctx context.Context, dtoID *entitydto.IDRequest) ([]orderprocessdto.OrderProcessDTO, error) {
	if processModels, err := s.r.GetProcessesByGroupItemID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.modelsToDTOs(processModels), nil
	}
}

func (s *OrderProcessService) modelsToDTOs(processModels []model.OrderProcess) []orderprocessdto.OrderProcessDTO {
	dtos := make([]orderprocessdto.OrderProcessDTO, 0)
	for _, processModel := range processModels {
		process := processModel.ToDomain()
		processDTO := &orderprocessdto.OrderProcessDTO{}
		processDTO.FromDomain(process)
		dtos = append(dtos, *processDTO)
	}

	return dtos
}
