package auth

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/validator"
)

type AuthUsecase struct {
	clientRepo domain.ClientRepository
	jwtSecret  string
}

func NewAuthUsecase(clientRepo domain.ClientRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		clientRepo: clientRepo,
		jwtSecret:  jwtSecret,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req RegisterDTO) (*AuthResponse, validator.ValidationErrors, error) {
	v := validator.New()
	v.CheckEmail("email", req.Email)
	v.CheckRequired("first_name", req.FirstName)
	v.CheckRequired("last_name", req.LastName)
	v.CheckMinLength("password", req.Password, 6)

	if !v.IsValid() {
		return nil, v, appErrors.ErrInvalidInput
	}

	// Check if email is already taken
	existing, err := u.clientRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		v.Add("email", "Email is already registered")
		return nil, v, appErrors.ErrDuplicateEntry
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	client := &domain.Client{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Company:      req.Company,
		Address1:     req.Address1,
		City:         req.City,
		Country:      req.Country,
		Phone:        req.Phone,
		Currency:     req.Currency,
		Status:       domain.ClientStatusActive,
	}

	if err := u.clientRepo.Create(ctx, client); err != nil {
		return nil, nil, err
	}

	token, err := auth.GenerateToken(u.jwtSecret, client.ID, client.Email, "client", 24*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	return &AuthResponse{
		Token: token,
		Client: ClientBrief{
			ID:        client.ID,
			Email:     client.Email,
			FirstName: client.FirstName,
			LastName:  client.LastName,
			Currency:  client.Currency,
		},
	}, nil, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req LoginDTO) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, appErrors.ErrInvalidInput
	}

	client, err := u.clientRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	if !auth.CheckPassword(req.Password, client.PasswordHash) {
		return nil, appErrors.ErrUnauthorized
	}

	if client.Status != domain.ClientStatusActive {
		return nil, errors.New("account is not active")
	}

	token, err := auth.GenerateToken(u.jwtSecret, client.ID, client.Email, "client", 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		Client: ClientBrief{
			ID:        client.ID,
			Email:     client.Email,
			FirstName: client.FirstName,
			LastName:  client.LastName,
			Currency:  client.Currency,
		},
	}, nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, clientID int64) (*ProfileResponse, error) {
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	balance, err := u.clientRepo.GetBalance(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &ProfileResponse{
		Client: ClientDetail{
			ID:        client.ID,
			Email:     client.Email,
			FirstName: client.FirstName,
			LastName:  client.LastName,
			Company:   client.Company,
			Address1:  client.Address1,
			Address2:  client.Address2,
			City:      client.City,
			State:     client.State,
			Postcode:  client.Postcode,
			Country:   client.Country,
			PhoneCC:   client.PhoneCC,
			Phone:     client.Phone,
			Currency:  client.Currency,
			Status:    string(client.Status),
		},
		Balance: balance,
	}, nil
}

func (u *AuthUsecase) UpdateProfile(ctx context.Context, clientID int64, req UpdateProfileDTO) (*ProfileResponse, error) {
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		client.FirstName = req.FirstName
	}
	if req.LastName != "" {
		client.LastName = req.LastName
	}
	client.Company = req.Company
	client.Address1 = req.Address1
	client.Address2 = req.Address2
	client.City = req.City
	client.State = req.State
	client.Postcode = req.Postcode
	if req.Country != "" {
		client.Country = req.Country
	}
	client.PhoneCC = req.PhoneCC
	client.Phone = req.Phone
	if req.Currency != "" {
		client.Currency = req.Currency
	}

	if err := u.clientRepo.Update(ctx, client); err != nil {
		return nil, err
	}

	return u.GetProfile(ctx, clientID)
}
