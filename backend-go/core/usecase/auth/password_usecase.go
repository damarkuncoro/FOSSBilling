package auth

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type PasswordUsecase struct {
	clientRepo domain.ClientRepository
}

func NewPasswordUsecase(clientRepo domain.ClientRepository) *PasswordUsecase {
	return &PasswordUsecase{clientRepo: clientRepo}
}

func (u *PasswordUsecase) ChangePassword(ctx context.Context, clientID int64, currentPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(currentPassword, client.PasswordHash) {
		return appErrors.ErrUnauthorized
	}

	hashed, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	client.PasswordHash = hashed
	return u.clientRepo.Update(ctx, client)
}
