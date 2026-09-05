package currency

import (
	"context"
	"errors"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)


var (
	ErrInvalidCurrencyCode = errors.New("currency code must be 3 characters (e.g. USD, IDR)")
	ErrCannotDeleteDefault = errors.New("cannot delete the default base currency")
)

type CreateCurrencyDTO struct {
	Code           string  `json:"code"`
	Title          string  `json:"title"`
	ConversionRate float64 `json:"conversion_rate"`
	Format         string  `json:"format"`
	PriceFormat    string  `json:"price_format"`
	IsDefault      bool    `json:"is_default"`
}

type UpdateCurrencyDTO struct {
	Title          string  `json:"title"`
	ConversionRate float64 `json:"conversion_rate"`
	Format         string  `json:"format"`
	PriceFormat    string  `json:"price_format"`
}

type CurrencyService struct {
	currencyRepo domain.CurrencyRepository
}

func NewCurrencyService(currencyRepo domain.CurrencyRepository) *CurrencyService {
	return &CurrencyService{currencyRepo: currencyRepo}
}

func (s *CurrencyService) ListCurrencies(ctx context.Context) ([]*domain.Currency, error) {
	return s.currencyRepo.List(ctx)
}

func (s *CurrencyService) GetCurrency(ctx context.Context, code string) (*domain.Currency, error) {
	return s.currencyRepo.GetByCode(ctx, code)
}

func (s *CurrencyService) CreateCurrency(ctx context.Context, dto CreateCurrencyDTO) (*domain.Currency, error) {
	code := strings.ToUpper(strings.TrimSpace(dto.Code))
	if len(code) != 3 {
		return nil, ErrInvalidCurrencyCode
	}

	if dto.ConversionRate <= 0 {
		dto.ConversionRate = 1.0
	}
	if dto.Format == "" {
		dto.Format = code + " {{price}}"
	}

	c := &domain.Currency{
		Code:           code,
		Title:          dto.Title,
		ConversionRate: dto.ConversionRate,
		Format:         dto.Format,
		PriceFormat:    dto.PriceFormat,
		IsDefault:      dto.IsDefault,
	}

	if err := s.currencyRepo.Create(ctx, c); err != nil {
		return nil, err
	}

	if dto.IsDefault {
		_ = s.currencyRepo.SetDefault(ctx, code)
	}

	return s.currencyRepo.GetByCode(ctx, code)
}

func (s *CurrencyService) UpdateCurrency(ctx context.Context, code string, dto UpdateCurrencyDTO) (*domain.Currency, error) {
	c, err := s.currencyRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if dto.Title != "" {
		c.Title = dto.Title
	}
	if dto.ConversionRate > 0 {
		c.ConversionRate = dto.ConversionRate
	}
	if dto.Format != "" {
		c.Format = dto.Format
	}
	if dto.PriceFormat != "" {
		c.PriceFormat = dto.PriceFormat
	}

	if err := s.currencyRepo.Update(ctx, c); err != nil {
		return nil, err
	}

	return s.currencyRepo.GetByCode(ctx, code)
}

func (s *CurrencyService) SetDefault(ctx context.Context, code string) error {
	return s.currencyRepo.SetDefault(ctx, code)
}

func (s *CurrencyService) DeleteCurrency(ctx context.Context, code string) error {
	c, err := s.currencyRepo.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if c.IsDefault {
		return ErrCannotDeleteDefault
	}
	return s.currencyRepo.Delete(ctx, code)
}
