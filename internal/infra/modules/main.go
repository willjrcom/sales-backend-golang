package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	emailservice "github.com/willjrcom/sales-backend-go/internal/infra/service/email"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func MainModules(db *bun.DB, chi *server.ServerChi, s3 *s3service.S3Client, rabbitmq *rabbitmq.RabbitMQ) {
	productRepository, _, _ := NewProductModule(db, chi)
	productCategoryRepository, _, _ := NewProductCategoryModule(db, chi)
	NewProductCategorySizeModule(db, chi)
	processRuleRepository, _, _ := NewProductCategoryProcessRuleModule(db, chi)

	addressRepository := NewAddressModule(db, chi)
	clientRepository, clientService, _ := NewClientModule(db, chi)
	employeeRepository, employeeService, _ := NewEmployeeModule(db, chi)
	contactRepository, _, _ := NewContactModule(db, chi)

	shiftRepository, shiftService, _ := NewShiftModule(db, chi)

	groupItemRepository, groupItemService, _ := NewGroupItemModule(db, chi)
	itemRepository, itemService, _ := NewItemModule(db, chi)

	orderProcessRepository, orderProcessService, _ := NewOrderProcessModule(db, chi)
	orderQueueRepository, orderQueueService, _ := NewOrderqueueModule(db, chi)

	// Stock module - deve ser inicializado antes do order module
	stockRepository, stockMovementRepository, _, stockService, _ := NewStockModule(db, chi)

	orderRepository, orderService, _ := NewOrderModule(db, chi)
	orderDeliveryRepository, orderDeliveryService, _ := NewOrderDeliveryModule(db, chi)
	deliveryDriverRepository, deliveryDriverService, _ := NewDeliveryDriverModule(db, chi)

	_, orderTableService, _ := NewOrderTableModule(db, chi)
	tableRepository, _, _ := NewTableModule(db, chi)
	NewPlaceModule(db, chi)

	_, orderPickupService, _ := NewOrderPickupModule(db, chi)

	// Usage Cost Repository (Creating here to pass to both Company and Fiscal modules)
	usageCostRepo := companyrepositorybun.NewCompanyUsageCostRepository(db)
	companySubscriptionRepo := companyrepositorybun.NewCompanySubscriptionRepositoryBun(db)

	companyRepository, companyService, checkoutUC, _ := NewCompanyModule(db, chi, usageCostRepo)
	sponsorRepository, _, _ := NewSponsorModule(db, chi)
	NewCompanyCategoryModule(db, chi)

	_, schemaService := NewSchemaModule(db, chi)
	userRepository, userService, _ := NewUserModule(db, chi)
	NewAdvertisingModule(db, chi, sponsorRepository, userRepository)

	emailService := emailservice.NewService(rabbitmq)

	go emailService.RunConsumer()
	// Fiscal invoice and usage cost modules
	_, _, _ = NewFiscalInvoiceModule(db, chi, companyRepository, companySubscriptionRepo, orderRepository, companyService, usageCostRepo)
	NewFiscalSettingsModule(db, chi, companyRepository, companyService)

	orderPrintService, _ := NewOrderPrintModule(db, chi)

	NewReportModule(db, chi)

	// Add S3 handler
	chi.AddHandler(handlerimpl.NewHandlerS3())
	// Public analytics handler for company/user listing
	chi.AddHandler(handlerimpl.NewHandlerPublicData(companyService, userService))

	checkoutUC.AddDependencies(userRepository)
	userService.AddDependencies(emailService)
	clientService.AddDependencies(contactRepository, companyService)
	employeeService.AddDependencies(contactRepository, userRepository, companyRepository)

	orderQueueService.AddDependencies(orderProcessRepository)
	orderProcessService.AddDependencies(orderQueueService, processRuleRepository, groupItemService, orderRepository, employeeService, groupItemRepository, orderService)

	itemService.AddDependencies(groupItemRepository, orderRepository, productRepository, productCategoryRepository, employeeRepository, orderService, groupItemService, stockRepository, stockMovementRepository)
	groupItemService.AddDependencies(itemRepository, productRepository, orderService, orderProcessService, employeeRepository, itemService)

	stockService.AddDependencies(productRepository, itemRepository, employeeRepository)

	orderService.AddDependencies(orderRepository, shiftRepository, productRepository, processRuleRepository, orderDeliveryRepository, stockRepository, stockMovementRepository, companySubscriptionRepo, groupItemService, orderProcessService, orderQueueService, orderDeliveryService, orderPickupService, orderTableService, companyService, employeeRepository, rabbitmq, clientService)
	orderDeliveryService.AddDependencies(addressRepository, clientRepository, orderRepository, orderService, deliveryDriverRepository, companyService, rabbitmq)
	deliveryDriverService.AddDependencies(employeeRepository)
	orderTableService.AddDependencies(tableRepository, orderService, companyService)
	orderPickupService.AddDependencies(orderService)

	shiftService.AddDependencies(employeeService, orderRepository, deliveryDriverRepository, orderProcessRepository, orderQueueRepository, processRuleRepository, employeeRepository)
	companyService.AddDependencies(addressRepository, *schemaService, userRepository, *userService, *employeeService, usageCostRepo, companySubscriptionRepo)

	orderPrintService.AddDependencies(orderService, orderRepository, shiftService, groupItemRepository, companyRepository, rabbitmq)
}
