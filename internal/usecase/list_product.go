package usecase

import "github.com/danyukod/go-kafka/internal/entity"

type ListProductOutputDto struct {
	ID    string
	Name  string
	Price float64
}

type ListProductUsecase struct {
	ProductRepository entity.ProductRepository
}

func NewListProductUsecase(productRepository entity.ProductRepository) *ListProductUsecase {
	return &ListProductUsecase{
		ProductRepository: productRepository,
	}
}

func (u *ListProductUsecase) Execute() ([]*ListProductOutputDto, error) {
	products, err := u.ProductRepository.FindAll()
	if err != nil {
		return nil, err
	}
	var output []*ListProductOutputDto
	for _, product := range products {
		output = append(output, &ListProductOutputDto{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		})
	}
	return output, nil
}
