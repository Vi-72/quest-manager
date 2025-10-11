package mock

import (
	"context"
	"errors"

	"quest-manager/internal/adapters/out/client/auth"

	"github.com/google/uuid"
)

// ConfigurableAuthClient allows configuring authentication behavior for different test scenarios.
type ConfigurableAuthClient struct {
	DefaultUserID uuid.UUID
	Behavior      AuthBehavior
}

// AuthBehavior определяет поведение mock auth client
type AuthBehavior int

const (
	// BehaviorSuccess - всегда успешная аутентификация
	BehaviorSuccess AuthBehavior = iota
	// BehaviorTokenExpired - токен протух
	BehaviorTokenExpired
	// BehaviorInvalidToken - невалидный токен
	BehaviorInvalidToken
	// BehaviorEmptyToken - пустой токен
	BehaviorEmptyToken
	// BehaviorMissingUser - отсутствует пользователь в ответе
	BehaviorMissingUser
)

// NewConfigurableAuthClient создает mock auth client с настраиваемым поведением
func NewConfigurableAuthClient(behavior AuthBehavior, userID uuid.UUID) *ConfigurableAuthClient {
	return &ConfigurableAuthClient{
		DefaultUserID: userID,
		Behavior:      behavior,
	}
}

// NewExpiredTokenAuthClient создает mock auth client который всегда возвращает ошибку протухшего токена
func NewExpiredTokenAuthClient() *ConfigurableAuthClient {
	return &ConfigurableAuthClient{
		Behavior: BehaviorTokenExpired,
	}
}

// NewInvalidTokenAuthClient создает mock auth client который всегда возвращает ошибку невалидного токена
func NewInvalidTokenAuthClient() *ConfigurableAuthClient {
	return &ConfigurableAuthClient{
		Behavior: BehaviorInvalidToken,
	}
}

// Authenticate реализует auth.Client интерфейс с настраиваемым поведением
func (a *ConfigurableAuthClient) Authenticate(ctx context.Context, jwtToken string) (uuid.UUID, error) {
	switch a.Behavior {
	case BehaviorSuccess:
		return a.DefaultUserID, nil

	case BehaviorTokenExpired:
		return uuid.Nil, auth.ErrTokenExpired

	case BehaviorInvalidToken:
		return uuid.Nil, errors.New("invalid jwt token format: token signature is invalid")

	case BehaviorEmptyToken:
		if jwtToken == "" {
			return uuid.Nil, errors.New("jwt token is empty")
		}
		return a.DefaultUserID, nil

	case BehaviorMissingUser:
		return uuid.Nil, errors.New("auth response missing user")

	default:
		return a.DefaultUserID, nil
	}
}
