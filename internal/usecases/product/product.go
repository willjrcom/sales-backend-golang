package productusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

var (
	ErrSizeIsInvalid = errors.New("size is invalid")
)

type Service struct {
	rp        model.ProductRepository
	rc        model.CategoryRepository
	s3Service *s3service.S3Client
}

func NewService(
	rp model.ProductRepository,
	rc model.CategoryRepository,
	s3 *s3service.S3Client,
) *Service {
	return &Service{
		rp:        rp,
		rc:        rc,
		s3Service: s3,
	}
}

func (s *Service) CreateProduct(ctx context.Context, dto *productcategorydto.ProductCreateDTO) (uuid.UUID, error) {
	product, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	if p, _ := s.rp.GetProductByCode(ctx, product.Code); p != nil {
		return uuid.Nil, errors.New("code product already exists")
	}

	if category, err := s.rc.GetCategoryById(ctx, product.CategoryID.String()); err == nil {
		product.Category = category
	} else {
		return uuid.Nil, err
	}

	if exists, err := product.FindSizeInCategory(); !exists {
		return uuid.Nil, err
	}

	if err := s.rp.CreateProduct(ctx, product); err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (s *Service) GetProductByCode(ctx context.Context, keys *productcategorydto.Keys) (*productcategorydto.ProductDTO, error) {
	if product, err := s.rp.GetProductByCode(ctx, keys.Code); err != nil {
		return nil, err
	} else {

		return productcategorydto.FromDomain(product), nil
	}
}

func (s *Service) UpdateProduct(ctx context.Context, dtoId *entitydto.IDRequest, dto *productcategorydto.ProductUpdateDTO) error {
	product, err := s.rp.GetProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateDomain(product); err != nil {
		return err
	}

	if p, _ := s.rp.GetProductByCode(ctx, product.Code); p != nil && p.ID != product.ID {
		return errors.New("code product already exists")
	}

	if err := s.rp.UpdateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProductById(ctx context.Context, dto *entitydto.IDRequest) error {
	product, err := s.rp.GetProductById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err := s.rp.DeleteProduct(ctx, dto.ID.String()); err != nil {
		return err
	}

	if product.ImagePath != nil {
		if err := s.s3Service.DeleteObject(*product.ImagePath); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetProductById(ctx context.Context, dto *entitydto.IDRequest) (*productcategorydto.ProductDTO, error) {
	if product, err := s.rp.GetProductById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return productcategorydto.FromDomain(product), nil
	}
}

func (s *Service) GetAllProducts(ctx context.Context) ([]productcategorydto.ProductDTO, error) {
	products, err := s.rp.GetAllProducts(ctx)

	if err != nil {
		return nil, err
	}

	dtos := productsToDtos(products)
	return dtos, nil
}

func productsToDtos(products []productentity.Product) []productcategorydto.ProductDTO {
	dtos := make([]productcategorydto.ProductDTO, len(products))
	for i, product := range products {
		dtos[i] = *productcategorydto.FromDomain(&product)
	}

	return dtos
}
