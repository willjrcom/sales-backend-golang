package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func MainModules(db *bun.DB, chi *server.ServerChi, s3 *s3service.S3Client) {
	productRepository, _, _ := NewProductModule(db, chi)
	productCategoryRepository, categoryService, _ := NewProductCategoryModule(db, chi)
	_, sizeService, _ := NewProductCategorySizeModule(db, chi)
	quantityRepository, quantityService, _ := NewProductCategoryQuantityModule(db, chi)
	processRuleRepository, _, _ := NewProductCategoryProcessRuleModule(db, chi)

	addressRepository := NewAddressModule(db, chi)
	clientRepository, clientService, _ := NewClientModule(db, chi)
	employeeRepository, employeeService, _ := NewEmployeeModule(db, chi)
	contactRepository, _, _ := NewContactModule(db, chi)

	shiftRepository, _, _ := NewShiftModule(db, chi)

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

	// Dependencies
	sizeService.AddDependencies(productCategoryRepository)
	quantityService.AddDependencies(productCategoryRepository)
	categoryService.AddDependencies(*quantityService, *sizeService)

	clientService.AddDependencies(contactRepository)
	employeeService.AddDependencies(contactRepository, userRepository)

	orderQueueService.AddDependencies(orderProcessRepository)
	orderProcessService.AddDependencies(orderQueueService, processRuleRepository, groupItemService, orderRepository, employeeService)

	itemService.AddDependencies(groupItemRepository, orderRepository, productRepository, quantityRepository)
	groupItemService.AddDependencies(itemRepository, productRepository)

	orderService.AddDependencies(shiftRepository, groupItemService, orderProcessService, processRuleRepository, orderQueueService)
	orderDeliveryService.AddDependencies(addressRepository, clientRepository, orderRepository, orderService, deliveryDriverRepository)
	deliveryDriverService.AddDependencies(employeeRepository)
	orderTableService.AddDependencies(tableRepository, orderService)
	orderPickupService.AddDependencies(orderService)

	companyService.AddDependencies(addressRepository, *schemaService, userRepository, *userService)
}
