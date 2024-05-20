package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/kafka"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func MainModules(db *bun.DB, chi *server.ServerChi, s3 *s3service.S3Client, producerKafka *kafka.KafkaProducer) {
	productRepository, productService, _ := NewProductCategoryProductModule(db, chi)
	productCategoryRepository, _, _ := NewProductCategoryModule(db, chi)
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
	NewDeliveryDriverModule(db, chi)

	_, orderTableService, _ := NewOrderTableModule(db, chi)
	tableRepository, _, _ := NewTableModule(db, chi)
	NewPlaceModule(db, chi)

	_, orderPickupService, _ := NewOrderPickupModule(db, chi)

	_, companyService, _ := NewCompanyModule(db, chi)
	_, schemaService := NewSchemaModule(db, chi)
	userRepository, userService, _ := NewUserModule(db, chi)

	// Dependencies
	productService.AddDependencies(productCategoryRepository, s3)
	sizeService.AddDependencies(productCategoryRepository)
	quantityService.AddDependencies(productCategoryRepository)

	clientService.AddDependencies(contactRepository)
	employeeService.AddDependencies(contactRepository)

	orderQueueService.AddDependencies(orderProcessRepository)
	orderProcessService.AddDependencies(orderQueueService, processRuleRepository, groupItemService, producerKafka)

	itemService.AddDependencies(groupItemRepository, orderRepository, productRepository, quantityRepository)
	groupItemService.AddDependencies(itemRepository, productRepository)

	orderService.AddDependencies(shiftRepository, groupItemService, orderProcessService, processRuleRepository, orderQueueService)
	orderDeliveryService.AddDependencies(addressRepository, clientRepository, orderRepository, employeeRepository, orderService)
	orderTableService.AddDependencies(tableRepository, orderService)
	orderPickupService.AddDependencies(orderService)

	companyService.AddDependencies(addressRepository, *schemaService, userRepository, *userService)
}
