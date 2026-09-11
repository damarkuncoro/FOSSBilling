package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
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
	issuer          string
	activityService *activity.ActivityService
}

func NewAuthUsecase(cr domain.ClientRepository, as *antispam.AntispamService, s string, issuer string, act *activity.ActivityService) *AuthUsecase {
	return &AuthUsecase{cr, as, s, issuer, act}
}

func (u *AuthUsecase) res(c *domain.Client) *AuthResponse {
	t, _ := auth.GenerateToken(u.jwtSecret, c.ID, c.Email, "client", 24*time.Hour)
	return &AuthResponse{Token: t, Client: ClientBrief{ID: c.ID, Email: c.Email, FirstName: c.FirstName, LastName: c.LastName, Currency: c.Currency}}
}

func (u *AuthUsecase) Register(ctx context.Context, req RegisterDTO, ip string) (*AuthResponse, validator.ValidationErrors, error) {
	v := validator.New(); v.CheckEmail("email", req.Email); v.CheckRequired("first_name", req.FirstName); v.CheckRequired("last_name", req.LastName); v.CheckPasswordStrength("password", req.Password, 8)
	if !v.IsValid() { return nil, v, appErrors.ErrInvalidInput }
	if u.antispamService != nil { if err := u.antispamService.ValidateSignup(ctx, req.Email, ip, req.Honeypot, req.CaptchaToken); err != nil { v.Add("email", err.Error()); return nil, v, appErrors.ErrInvalidInput } }
	if ex, _ := u.clientRepo.GetByEmail(ctx, req.Email); ex != nil { return nil, nil, appErrors.ErrDuplicate }
	hp, _ := auth.HashPassword(req.Password); c := &domain.Client{Email: req.Email, PasswordHash: hp, FirstName: req.FirstName, LastName: req.LastName, Company: req.Company, Address1: req.Address1, City: req.City, Country: req.Country, Phone: req.Phone, Currency: req.Currency, ReferrerID: req.ReferrerID, Status: domain.ClientStatusActive}
	if err := u.clientRepo.Create(ctx, c); err != nil { return nil, nil, err }
	return u.res(c), nil, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req LoginDTO) (*AuthResponse, error) {
	c, err := u.clientRepo.GetByEmail(ctx, req.Email); if err != nil || !auth.CheckPassword(req.Password, c.PasswordHash) { return nil, appErrors.ErrUnauthorized }
	if c.Status != domain.ClientStatusActive { return nil, errors.New("inactive") }
	if c.TwoFactorEnabled { return &AuthResponse{TwoFactorRequired: true, Client: ClientBrief{ID: c.ID, Email: c.Email, FirstName: c.FirstName}}, nil }

	if u.activityService != nil {
		_ = u.activityService.LogClientEvent(ctx, c.ID, "login", "Client logged into portal", "")
	}

	return u.res(c), nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, id int64) (*ProfileResponse, error) {
	c, err := u.clientRepo.GetByID(ctx, id); if err != nil { return nil, err }
	b, _ := u.clientRepo.GetBalance(ctx, id); bd := ""; if c.Birthday != nil { bd = c.Birthday.Format("2006-01-02") }
	cd := ClientDetail{ID: c.ID, AID: c.AID, Email: c.Email, FirstName: c.FirstName, LastName: c.LastName, Gender: c.Gender, Birthday: bd, Company: c.Company, CompanyVat: c.CompanyVat, CompanyNumber: c.CompanyNumber, Type: c.Type, Address1: c.Address1, Address2: c.Address2, City: c.City, State: c.State, Postcode: c.Postcode, Country: c.Country, PhoneCC: c.PhoneCC, Phone: c.Phone, Currency: c.Currency, BillingEmail: c.BillingEmail, Status: string(c.Status)}
	cd.Custom1, cd.Custom2, cd.Custom3, cd.Custom4, cd.Custom5, cd.Custom6, cd.Custom7, cd.Custom8, cd.Custom9, cd.Custom10 = c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5, c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10
	cd.Custom11, cd.Custom12, cd.Custom13, cd.Custom14, cd.Custom15, cd.Custom16, cd.Custom17, cd.Custom18, cd.Custom19, cd.Custom20 = c.Custom11, c.Custom12, c.Custom13, c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20
	return &ProfileResponse{Client: cd, Balance: b}, nil
}

func (u *AuthUsecase) UpdateProfile(ctx context.Context, id int64, req UpdateProfileDTO) (*ProfileResponse, error) {
	c, err := u.clientRepo.GetByID(ctx, id); if err != nil { return nil, err }
	if req.FirstName != "" { c.FirstName = security.SanitizeAlphaNumeric(req.FirstName) }
	if req.LastName != "" { c.LastName = security.SanitizeAlphaNumeric(req.LastName) }
	c.Gender = security.SanitizeAlphaNumeric(req.Gender); if t, err := time.Parse("2006-01-02", req.Birthday); err == nil { c.Birthday = &t }
	c.Company, c.CompanyVat, c.CompanyNumber, c.Type = security.SanitizeHTML(req.Company), req.CompanyVat, req.CompanyNumber, req.Type
	c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country = req.Address1, req.Address2, req.City, req.State, req.Postcode, req.Country
	c.PhoneCC, c.Phone, c.Currency, c.BillingEmail = req.PhoneCC, req.Phone, req.Currency, req.BillingEmail
	c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5, c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10 = req.Custom1, req.Custom2, req.Custom3, req.Custom4, req.Custom5, req.Custom6, req.Custom7, req.Custom8, req.Custom9, req.Custom10
	c.Custom11, c.Custom12, c.Custom13, c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20 = req.Custom11, req.Custom12, req.Custom13, req.Custom14, req.Custom15, req.Custom16, req.Custom17, req.Custom18, req.Custom19, req.Custom20
	if err := u.clientRepo.Update(ctx, c); err != nil { return nil, err }
	return u.GetProfile(ctx, id)
}

func (u *AuthUsecase) VerifyTwoFactor(ctx context.Context, em, co string) (*AuthResponse, error) {
	c, err := u.clientRepo.GetByEmail(ctx, em); if err != nil || c.TwoFactorSecret == nil || !security.VerifyTOTP(*c.TwoFactorSecret, co) { return nil, appErrors.ErrUnauthorized }
	return u.res(c), nil
}

func (u *AuthUsecase) SetupTwoFactor(ctx context.Context, id int64) (*TwoFactorSetupResponse, error) {
	c, _ := u.clientRepo.GetByID(ctx, id); s := security.GenerateTOTPSecret(); c.TwoFactorSecret = &s; _ = u.clientRepo.Update(ctx, c)
	return &TwoFactorSetupResponse{Secret: s, QRURL: security.GenerateTOTPURL(c.Email, u.issuer, s)}, nil
}

func (u *AuthUsecase) EnableTwoFactor(ctx context.Context, id int64, co string) error {
	c, _ := u.clientRepo.GetByID(ctx, id); if c.TwoFactorSecret == nil || !security.VerifyTOTP(*c.TwoFactorSecret, co) { return errors.New("invalid") }
	c.TwoFactorEnabled = true; return u.clientRepo.Update(ctx, c)
}

func (u *AuthUsecase) DisableTwoFactor(ctx context.Context, id int64) error {
	c, _ := u.clientRepo.GetByID(ctx, id); c.TwoFactorEnabled, c.TwoFactorSecret = false, nil; return u.clientRepo.Update(ctx, c)
}

func (u *AuthUsecase) AdminImpersonateClient(ctx context.Context, id int64) (string, error) {
	c, err := u.clientRepo.GetByID(ctx, id); if err != nil { return "", err }
	return auth.GenerateTokenExt(u.jwtSecret, c.ID, c.Email, "client", time.Hour, true)
}

func (u *AuthUsecase) OAuthLoginOrRegister(ctx context.Context, provider, oauthID, email, name string) (*AuthResponse, error) {
	c, err := u.clientRepo.GetByOAuth(ctx, provider, oauthID)
	if err == nil {
		if c.Status != domain.ClientStatusActive { return nil, errors.New("inactive") }
		return u.res(c), nil
	}

	// Try by email
	c, err = u.clientRepo.GetByEmail(ctx, email)
	if err == nil {
		// Link account
		c.OAuthProvider = provider
		c.OAuthID = oauthID
		_ = u.clientRepo.Update(ctx, c)
		return u.res(c), nil
	}

	// Register new
	names := strings.SplitN(name, " ", 2)
	fname := names[0]; lname := ""; if len(names) > 1 { lname = names[1] }
	hp, _ := auth.HashPassword(security.GenerateTOTPSecret()) // Random password for OAuth users
	c = &domain.Client{
		Email: email, PasswordHash: hp, FirstName: fname, LastName: lname,
		OAuthProvider: provider, OAuthID: oauthID, Status: domain.ClientStatusActive,
	}
	if err := u.clientRepo.Create(ctx, c); err != nil { return nil, err }
	return u.res(c), nil
}
