package api

import (
	"context"
	"fmt"
)

// UserService предоставляет методы для работы с пользователями
type UserService struct {
	client *Client
}

// NewUserService создает новый сервис для работы с пользователями
func NewUserService(client *Client) *UserService {
	return &UserService{client: client}
}

// Get возвращает информацию о пользователе по ID
func (s *UserService) Get(ctx context.Context, customerID, userID string) (*User, error) {
	var result User
	path := fmt.Sprintf("/api/v1/customers/%s/users/%s", customerID, userID)
	err := s.client.Get(ctx, path, nil, &result)
	return &result, err
}

// GetByEmail возвращает информацию о пользователе по email
func (s *UserService) GetByEmail(ctx context.Context, customerID, email string) (*User, error) {
	var result User
	query := map[string]string{
		"email": email,
	}
	path := fmt.Sprintf("/api/v1/customers/%s/users", customerID)
	err := s.client.Get(ctx, path, query, &result)
	return &result, err
}
