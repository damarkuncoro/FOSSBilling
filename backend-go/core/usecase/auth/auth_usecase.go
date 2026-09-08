package auth

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/validator"
)

type AuthUsecase struct {
	clientRepo      domain.ClientRepository
	antispamService *antispam.AntispamService
	jwtSecret       string
}

func NewAuthUsecase(clientRepo domain.ClientRepository, antispamService *antispam.AntispamService, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		clientRepo:      clientRepo,
		antispamService: antispamService,
		jwtSecret:       jwtSecret,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req RegisterDTO, remoteIP string) (*AuthResponse, validator.ValidationErrors, error) {
	v := validator.New()
	v.CheckEmail("email", req.Email)
	v.CheckRequired("first_name", req.FirstName)
	v.CheckRequired("last_name", req.LastName)
	v.CheckMinLength("password", req.Password, 6)

	if !v.IsValid() {
		return nil, v, appErrors.ErrInvalidInput
	}

	// Anti-spam validation
	if u.antispamService != nil {
		if err := u.antispamService.ValidateSignup(ctx, req.Email, remoteIP, req.Honeypot, req.CaptchaToken); err != nil {
			v.Add("email", err.Error())
			return nil, v, appErrors.ErrInvalidInput
		}
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

	if client.TwoFactorEnabled {
		return &AuthResponse{
			TwoFactorRequired: true,
			Client: ClientBrief{
				ID:        client.ID,
				Email:     client.Email,
				FirstName: client.FirstName,
			},
		}, nil
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

func (u *AuthUsecase) VerifyTwoFactor(ctx context.Context, email, code string) (*AuthResponse, error) {
	client, err := u.clientRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	if client.TwoFactorSecret == nil || !security.VerifyTOTP(*client.TwoFactorSecret, code) {
		return nil, errors.New("invalid two-factor code")
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

func (u *AuthUsecase) SetupTwoFactor(ctx context.Context, clientID int64) (*TwoFactorSetupResponse, error) {
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	secret := security.GenerateTOTPSecret()
	qrURL := security.GenerateTOTPURL(client.Email, "FOSSBilling", secret)

	// Save secret temporarily but don't enable yet
	client.TwoFactorSecret = &secret
	if err := u.clientRepo.Update(ctx, client); err != nil {
		return nil, err
	}

	return &TwoFactorSetupResponse{
		Secret: secret,
		QRURL:  qrURL,
	}, nil
}

func (u *AuthUsecase) EnableTwoFactor(ctx context.Context, clientID int64, code string) error {
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if client.TwoFactorSecret == nil {
		return errors.New("2FA secret not generated")
	}

	if !security.VerifyTOTP(*client.TwoFactorSecret, code) {
		return errors.New("invalid verification code")
	}

	client.TwoFactorEnabled = true
	return u.clientRepo.Update(ctx, client)
}

func (u *AuthUsecase) DisableTwoFactor(ctx context.Context, clientID int64) error {
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	client.TwoFactorEnabled = false
	client.TwoFactorSecret = nil
	return u.clientRepo.Update(ctx, client)
}
