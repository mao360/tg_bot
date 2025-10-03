package payment

import (
	"context"
	"tg_bot/internal/models"
	"tg_bot/pkg/logger"
)

type Service struct {
	logger logger.Logger
	products map[string]models.Product
}

func NewService(logger logger.Logger) *Service {
	// Инициализируем продукты
	products := map[string]models.Product{
		"premium": {
			ID:          "premium",
			Name:        "Премиум подписка",
			Description: "Доступ к премиум функциям на месяц",
			Price:       29900, // 299 рублей
			Currency:    "RUB",
		},
		"pro": {
			ID:          "pro",
			Name:        "Pro подписка",
			Description: "Расширенные возможности на 3 месяца",
			Price:       79900, // 799 рублей
			Currency:    "RUB",
		},
	}
	
	return &Service{
		logger: logger,
		products: products,
	}
}

func (s *Service) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	product, exists := s.products[productID]
	if !exists {
		return nil, ErrProductNotFound
	}
	return &product, nil
}

func (s *Service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	for _, product := range s.products {
		products = append(products, product)
	}
	return products, nil
}

func (s *Service) ProcessPayment(ctx context.Context, payment *models.Payment) error {
	s.logger.Info("Processing payment", 
		"user_id", payment.UserID,
		"amount", payment.Amount,
		"currency", payment.Currency,
		"product_id", payment.ProductID)
	
	// Здесь будет логика обработки платежа
	// В будущем здесь будет интеграция с платежными системами
	
	return nil
}

func (s *Service) ValidatePreCheckout(ctx context.Context, query *models.PreCheckoutQuery) error {
	s.logger.Debug("Validating pre-checkout", 
		"amount", query.TotalAmount,
		"currency", query.Currency,
		"payload", query.InvoicePayload)
	
	// Здесь будет валидация платежа
	// Проверка суммы, валюты, payload и т.д.
	
	return nil
}
