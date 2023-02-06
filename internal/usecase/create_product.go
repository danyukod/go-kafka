package usecase

import "github.com/danyukod/go-kafka/internal/entity"

type CreateProductInputDto struct {
	Name  string
	Price float64
}

type CreateProductOutputDto struct {
	ID    string
	Name  string
	Price float64
}

type CreateProductUsecase struct {
	ProductRepository entity.ProductRepository
}

func (u *CreateProductUsecase) Execute(input CreateProductInputDto) (*CreateProductOutputDto, error) {
	product := entity.NewProduct(input.Name, input.Price)
	if err := u.ProductRepository.Create(product); err != nil {
		return nil, err
	}
	return &CreateProductOutputDto{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
