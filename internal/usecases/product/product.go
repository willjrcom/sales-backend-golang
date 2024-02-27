package productusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

var (
	ErrSizeIsInvalid = errors.New("size is invalid")
)

type Service struct {
	rp productentity.ProductRepository
	rc productentity.CategoryRepository
}

func NewService(r productentity.ProductRepository, c productentity.CategoryRepository) *Service {
	return &Service{rp: r, rc: c}
}

func (s *Service) RegisterProduct(ctx context.Context, dto *productdto.RegisterProductInput) (uuid.UUID, error) {
	product, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if category, err := s.rc.GetCategoryById(ctx, product.CategoryID.String()); err == nil {
		product.Category = category
	} else {
		return uuid.Nil, err
	}

	if exists, err := product.FindSizeInCategory(); !exists {
		return uuid.Nil, err
	}

	// imagePath, err := s3.UploadToS3(dto.Image)
	// if err != nil {
	// 	fmt.Printf("Erro ao fazer upload: %v\n", err)
	// }

	// product.ImagePath = imagePath

	if err := s.rp.RegisterProduct(ctx, product); err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (s *Service) UpdateProduct(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateProductInput) error {
	product, err := s.rp.GetProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(product); err != nil {
		return err
	}

	if err := s.rp.UpdateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProductById(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rp.GetProductById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rp.DeleteProduct(ctx, dto.ID.String()); err != nil {
		return err
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
