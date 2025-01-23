// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package sizeusecases

import (
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size"
)

// Injectors from wire.go:

func InitializeService() (*Service, error) {
	db := database.NewPostgreSQLConnection()
	sizeRepository := sizerepositorybun.NewSizeRepositoryBun(db)
	categoryRepository := productcategoryrepositorybun.NewProductCategoryRepositoryBun(db)
	service := NewService(sizeRepository, categoryRepository)
	return service, nil
}

func InitializeServiceForTest() (*Service, error) {
	sizeRepository := sizerepositorylocal.NewSizeRepositoryLocal()
	categoryRepository := categoryrepositorylocal.NewCategoryRepositoryLocal()
	service := NewService(sizeRepository, categoryRepository)
	return service, nil
}
