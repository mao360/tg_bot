package user

import (
	"context"
	"tg_bot/internal/models"
	"tg_bot/pkg/logger"
)

type Service struct {
	logger logger.Logger
	states map[int64]string // Временное хранение состояний пользователей
}

func NewService(logger logger.Logger) *Service {
	return &Service{
		logger: logger,
		states: make(map[int64]string),
	}
}

func (s *Service) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*models.User, error) {
	// Пока создаем пользователя в памяти
	// В будущем здесь будет работа с базой данных
	user := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Language:   "ru",
		IsPremium:  false,
	}
	
	s.logger.Debug("User created/retrieved", "telegram_id", telegramID, "username", username)
	return user, nil
}

func (s *Service) SetUserState(ctx context.Context, userID int64, state string) {
	s.states[userID] = state
	s.logger.Debug("User state set", "user_id", userID, "state", state)
}

func (s *Service) GetUserState(ctx context.Context, userID int64) (string, bool) {
	state, exists := s.states[userID]
	return state, exists
}

func (s *Service) ActivatePremium(ctx context.Context, userID int64) error {
	s.logger.Info("Premium activated for user", "user_id", userID)
	// Здесь будет логика активации премиум функций
	return nil
}
