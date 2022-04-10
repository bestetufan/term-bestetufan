package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/domain/repo"
)

type StoreService struct {
	categoryRepo repo.CategoryRepository
	productRepo  repo.ProductRepository
	basketRepo   repo.BasketRepository
	orderRepo    repo.OrderRepository
}

func NewStoreService(cr repo.CategoryRepository, pr repo.ProductRepository, br repo.BasketRepository, or repo.OrderRepository) *StoreService {
	return &StoreService{
		categoryRepo: cr,
		productRepo:  pr,
		basketRepo:   br,
		orderRepo:    or,
	}
}

func (s *StoreService) GetAllCategories(pageIndex, pageSize int, onlyActives bool) ([]entity.Category, int) {
	if onlyActives {
		return s.categoryRepo.GetAllActives(pageIndex, pageSize)
	}
	return s.categoryRepo.GetAll(pageIndex, pageSize)
}

func (s *StoreService) GetCategory(categoryId uint32) *entity.Category {
	return s.categoryRepo.GetById(categoryId)
}

func (s *StoreService) CreateCategory(category *entity.Category) error {
	categoryExists := s.categoryRepo.GetByName(category.Name)
	if categoryExists != nil {
		return errors.New("category with same name already exist in database")
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return errors.New("an unknown error occurred during operation")
	}

	return nil
}

func (s *StoreService) CreateBulkCategory(file multipart.File) (int, int, error) {
	reader := csv.NewReader(file)

	reader.Comma = ';'
	lines, err := reader.ReadAll()
	if err != nil {
		return 0, 0, errors.New("unable to initialize csv reader")
	}

	addedCount := 0
	existingCount := 0

	for _, line := range lines[1:] {
		name := line[0]
		isActive, _ := strconv.ParseBool(line[1])

		category := entity.NewCategory(
			name,     // Name
			isActive, // IsActive
		)

		categoryExists := s.categoryRepo.GetByName(name)
		if categoryExists != nil {
			existingCount += 1
		} else {
			if err := s.categoryRepo.Create(category); err != nil {
				return addedCount, existingCount, fmt.Errorf("unable to create category: %s", name)
			}
			addedCount += 1
		}
	}

	return addedCount, existingCount, nil
}

func (s *StoreService) GetAllProducts(pageIndex, pageSize int) ([]entity.Product, int) {
	items, count := s.productRepo.GetAll(pageIndex, pageSize)

	return items, count
}

func (s *StoreService) GetProduct(productId uint32) *entity.Product {
	return s.productRepo.GetById(productId)
}

func (s *StoreService) CreateProduct(product *entity.Product) error {
	productExists := s.productRepo.GetBySKU(product.Sku)
	if productExists != nil {
		return errors.New("product with same sku already exist in database")
	}

	err := s.productRepo.Create(product)
	if err != nil {
		return errors.New("an unknown error occurred during operation")
	}

	return nil
}

func (s *StoreService) GetBasket(userName string) *entity.Basket {
	return s.basketRepo.Get(userName)
}

func (s *StoreService) AddItemToBasket(userName string, productId uint32, quantity int) error {
	basket := s.basketRepo.Get(userName)
	if basket == nil {
		basket, _ := entity.NewBasket(userName)
		if err := s.basketRepo.Create(basket); err != nil {
			return errors.New("unable to create basket")
		}
	}

	product := s.productRepo.GetById(productId)
	if product == nil {
		return errors.New("product not found")
	}

	if product.Quantity < quantity {
		return errors.New("not enough stock")
	}

	item, _ := entity.NewBasketItem(basket.ID, productId, quantity)
	if err := s.basketRepo.AddItem(item); err != nil {
		return errors.New("unable to add item to basket")
	}

	product.Quantity -= quantity
	if err := s.productRepo.Update(product); err != nil {
		return errors.New("unable to update product stock info")
	}

	return nil
}

func (s *StoreService) UpdateItemInBasket(userName string, productId uint32, quantity int) error {
	basket := s.basketRepo.Get(userName)
	if basket == nil {
		return errors.New("basket not found")
	}

	_, item := basket.SearchItemByProductId(productId)
	if item == nil {
		return errors.New("item not found")
	}

	product := s.productRepo.GetById(productId)
	if product == nil {
		return errors.New("product not found")
	}

	if product.Quantity < quantity {
		return errors.New("not enough stock")
	}

	product.Quantity = product.Quantity - (quantity - item.Quantity)
	if err := s.productRepo.Update(product); err != nil {
		return errors.New("unable to update product stock info")
	}

	item.Quantity = quantity
	if err := s.basketRepo.UpdateItem(item); err != nil {
		return errors.New("unable to update item info")
	}

	return nil
}

func (s *StoreService) RemoveItemFromBasket(userName string, productId uint32) error {
	basket := s.basketRepo.Get(userName)
	if basket == nil {
		return errors.New("basket not found")
	}

	_, item := basket.SearchItemByProductId(productId)
	if item == nil {
		return errors.New("item not found")
	}

	product := s.productRepo.GetById(productId)
	if product == nil {
		return errors.New("product not found")
	}

	product.Quantity += item.Quantity
	if err := s.productRepo.Update(product); err != nil {
		return errors.New("unable to update product stock info")
	}

	if err := s.basketRepo.DeleteItem(item); err != nil {
		return errors.New("unable to hard remove item")
	}

	return nil
}

func (s *StoreService) GetAllOrders(userName string) []entity.Order {
	items := s.orderRepo.GetAll(userName)
	return items
}

func (s *StoreService) CreateOrder(userName string, name string, address string, phoneNumber string,
	cardNumber string, cardExp string, cardCVV int) error {
	order, _ := entity.NewOrder(userName, name, address, phoneNumber, cardNumber, cardExp, cardCVV)

	basket := s.basketRepo.Get(userName)
	if basket == nil {
		return errors.New("basket not found")
	}

	for _, v := range basket.Items {
		item, _ := entity.NewOrderItem("", v.ProductID)
		if err := order.AddItem(item); err != nil {
			return errors.New("unable to add order item")
		}
	}

	if err := s.orderRepo.Create(order); err != nil {
		return errors.New("unable to create order")
	}

	if err := s.basketRepo.DeleteItemsByBasketId(basket.ID); err != nil {
		return errors.New("unable to clear basket")
	}

	return nil
}

func (s *StoreService) CancelOrder(userName string, orderId string) error {
	order := s.orderRepo.Get(userName, orderId)
	if order == nil {
		return errors.New("order not found")
	}

	now := time.Now()
	orderDate := order.CreatedAt.In(time.UTC)
	days := now.Sub(orderDate).Hours() / 24
	if days > 14 {
		return errors.New("order cancel period ended")
	}

	order.Status = "canceled"
	if err := s.orderRepo.Update(order); err != nil {
		return errors.New("unable to update order status")
	}

	return nil
}
