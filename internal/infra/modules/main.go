package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func MainModules(db *bun.DB, chi *server.ServerChi, s3 *s3service.S3Client) {
	productRepository, _, _ := NewProductModule(db, chi)
	productCategoryRepository, _, _ := NewProductCategoryModule(db, chi)
	NewProductCategorySizeModule(db, chi)
	quantityRepository, _, _ := NewProductCategoryQuantityModule(db, chi)
	processRuleRepository, _, _ := NewProductCategoryProcessRuleModule(db, chi)

	addressRepository := NewAddressModule(db, chi)
	clientRepository, clientService, _ := NewClientModule(db, chi)
	employeeRepository, employeeService, _ := NewEmployeeModule(db, chi)
	contactRepository, _, _ := NewContactModule(db, chi)

	shiftRepository, shiftService, _ := NewShiftModule(db, chi)

	groupItemRepository, groupItemService, _ := NewGroupItemModule(db, chi)
	itemRepository, itemService, _ := NewItemModule(db, chi)

	orderProcessRepository, orderProcessService, _ := NewOrderProcessModule(db, chi)
	_, orderQueueService, _ := NewOrderqueueModule(db, chi)

	orderRepository, orderService, _ := NewOrderModule(db, chi)
	_, orderDeliveryService, _ := NewOrderDeliveryModule(db, chi)
	deliveryDriverRepository, deliveryDriverService, _ := NewDeliveryDriverModule(db, chi)

	_, orderTableService, _ := NewOrderTableModule(db, chi)
	tableRepository, _, _ := NewTableModule(db, chi)
	NewPlaceModule(db, chi)

	_, orderPickupService, _ := NewOrderPickupModule(db, chi)

	_, companyService, _ := NewCompanyModule(db, chi)
	_, schemaService := NewSchemaModule(db, chi)
	userRepository, userService, _ := NewUserModule(db, chi)

	orderPrintService, _ := NewOrderPrintModule(db, chi)

	NewReportModule(db, chi)

	clientService.AddDependencies(contactRepository)
	employeeService.AddDependencies(contactRepository, userRepository)

	orderQueueService.AddDependencies(orderProcessRepository)
	orderProcessService.AddDependencies(orderQueueService, processRuleRepository, groupItemService, orderRepository, employeeService, groupItemRepository)

	itemService.AddDependencies(groupItemRepository, orderRepository, productRepository, quantityRepository, productCategoryRepository, orderService, groupItemService)
	groupItemService.AddDependencies(itemRepository, productRepository, orderService, orderProcessService)

	orderService.AddDependencies(orderRepository, shiftRepository, productRepository, processRuleRepository, groupItemService, orderProcessService, orderQueueService, orderDeliveryService, orderPickupService, orderTableService)
	orderDeliveryService.AddDependencies(addressRepository, clientRepository, orderRepository, orderService, deliveryDriverRepository)
	deliveryDriverService.AddDependencies(employeeRepository)
	orderTableService.AddDependencies(tableRepository, orderService, companyService)
	orderPickupService.AddDependencies(orderService)

	shiftService.AddDependencies(employeeService, orderRepository)
	companyService.AddDependencies(addressRepository, *schemaService, userRepository, *userService)

	orderPrintService.AddDependencies(orderService, orderRepository, shiftService, groupItemRepository)
}
