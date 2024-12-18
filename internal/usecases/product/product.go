package productusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

var (
	ErrSizeIsInvalid = errors.New("size is invalid")
)

type Service struct {
	rp        productentity.ProductRepository
	rc        productentity.CategoryRepository
	s3Service *s3service.S3Client
}

func NewService(
	rp *productrepositorybun.ProductRepositoryBun,
	rc *productcategoryrepositorybun.ProductCategoryRepositoryBun,
	s3 *s3service.S3Client,
) *Service {
	return &Service{
		rp:        rp,
		rc:        rc,
		s3Service: s3,
	}
}

func (s *Service) CreateProduct(ctx context.Context, dto *productdto.CreateProductInput) (uuid.UUID, error) {
	product, err := dto.ToModel()

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

func (s *Service) GetProductByCode(ctx context.Context, keys *productdto.Keys) (*productdto.ProductOutput, error) {
	if product, err := s.rp.GetProductByCode(ctx, keys.Code); err != nil {
		return nil, err
	} else {
		dtoOutput := &productdto.ProductOutput{}
		dtoOutput.FromModel(product)
		return dtoOutput, nil
	}
}

func (s *Service) UpdateProduct(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateProductInput) error {
	product, err := s.rp.GetProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(product); err != nil {
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

func (s *Service) DeleteProductById(ctx context.Context, dto *entitydto.IdRequest) error {
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

func (s *Service) GetProductById(ctx context.Context, dto *entitydto.IdRequest) (*productdto.ProductOutput, error) {
	if product, err := s.rp.GetProductById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dtoOutput := &productdto.ProductOutput{}
		dtoOutput.FromModel(product)
		return dtoOutput, nil
	}
}

func (s *Service) GetAllProducts(ctx context.Context) ([]productdto.ProductOutput, error) {
	products, err := s.rp.GetAllProducts(ctx)

	if err != nil {
		return nil, err
	}

	dtos := productsToDtos(products)
	return dtos, nil
}

func productsToDtos(products []productentity.Product) []productdto.ProductOutput {
	dtos := make([]productdto.ProductOutput, len(products))
	for i, product := range products {
		dtos[i].FromModel(&product)
	}

	return dtos
}
