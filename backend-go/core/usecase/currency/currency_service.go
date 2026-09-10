package currency

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrInvalidCode = errors.New("invalid code (must be 3 chars)")
	ErrDeleteDef   = errors.New("cannot delete default currency")
)

type CreateCurrencyDTO struct { Code, Title, Format, PriceFormat string; ConversionRate float64; IsDefault bool }
type UpdateCurrencyDTO struct { Title, Format, PriceFormat string; ConversionRate float64 }

type CurrencyService struct{ repo domain.CurrencyRepository }

func NewCurrencyService(r domain.CurrencyRepository) *CurrencyService { return &CurrencyService{r} }

func (s *CurrencyService) ListCurrencies(ctx context.Context) ([]*domain.Currency, error) { return s.repo.List(ctx) }
func (s *CurrencyService) GetCurrency(ctx context.Context, c string) (*domain.Currency, error) { return s.repo.GetByCode(ctx, c) }

func (s *CurrencyService) CreateCurrency(ctx context.Context, d CreateCurrencyDTO) (*domain.Currency, error) {
	c := strings.ToUpper(strings.TrimSpace(d.Code)); if len(c) != 3 { return nil, ErrInvalidCode }
	if d.ConversionRate <= 0 { d.ConversionRate = 1.0 }; if d.Format == "" { d.Format = c + " {{price}}" }
	if err := s.repo.Create(ctx, &domain.Currency{Code: c, Title: d.Title, ConversionRate: d.ConversionRate, Format: d.Format, PriceFormat: d.PriceFormat, IsDefault: d.IsDefault}); err != nil { return nil, err }
	if d.IsDefault { _ = s.repo.SetDefault(ctx, c) }; return s.repo.GetByCode(ctx, c)
}

func (s *CurrencyService) UpdateCurrency(ctx context.Context, c string, d UpdateCurrencyDTO) (*domain.Currency, error) {
	curr, err := s.repo.GetByCode(ctx, c); if err != nil { return nil, err }
	if d.Title != "" { curr.Title = d.Title }; if d.ConversionRate > 0 { curr.ConversionRate = d.ConversionRate }
	if d.Format != "" { curr.Format = d.Format }; if d.PriceFormat != "" { curr.PriceFormat = d.PriceFormat }
	if err := s.repo.Update(ctx, curr); err != nil { return nil, err }; return s.repo.GetByCode(ctx, c)
}

func (s *CurrencyService) SetDefault(ctx context.Context, c string) error { return s.repo.SetDefault(ctx, c) }

func (s *CurrencyService) DeleteCurrency(ctx context.Context, c string) error {
	curr, err := s.repo.GetByCode(ctx, c); if err != nil { return err }
	if curr.IsDefault { return ErrDeleteDef }; return s.repo.Delete(ctx, c)
}

func (s *CurrencyService) UpdateExchangeRates(ctx context.Context) error {
	cs, err := s.repo.List(ctx); if err != nil { return err }
	base := "USD"; for _, c := range cs { if c.IsDefault { base = c.Code; break } }
	resp, err := http.Get("https://open.er-api.com/v6/latest/" + base); if err != nil { return err }; defer resp.Body.Close()
	var d struct{ Rates map[string]float64 `json:"rates"` }; if err := json.NewDecoder(resp.Body).Decode(&d); err != nil { return err }
	for _, c := range cs { if !c.IsDefault { if r, ok := d.Rates[c.Code]; ok { c.ConversionRate = r; _ = s.repo.Update(ctx, c) } } }
	return nil
}
