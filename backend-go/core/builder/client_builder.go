package builder

import (
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

type ClientBuilder struct { em, pw, fn, ln, co, cnt, cur string; st domain.ClientStatus }

func NewClientBuilder() *ClientBuilder { return &ClientBuilder{cnt: "ID", cur: "USD", st: domain.ClientStatusActive} }
func (b *ClientBuilder) WithEmail(e string) *ClientBuilder { b.em = e; return b }
func (b *ClientBuilder) WithPassword(p string) *ClientBuilder { b.pw = p; return b }
func (b *ClientBuilder) WithName(f, l string) *ClientBuilder { b.fn, b.ln = f, l; return b }
func (b *ClientBuilder) WithCompany(c string) *ClientBuilder { b.co = c; return b }
func (b *ClientBuilder) WithCountryAndCurrency(cnt, cur string) *ClientBuilder { b.cnt, b.cur = cnt, cur; return b }

func (b *ClientBuilder) Build() (*domain.Client, error) {
	if b.em == "" { return nil, errors.New("email required") }
	if b.pw == "" { b.pw = "Password123!" }; hp, _ := auth.HashPassword(b.pw); n := time.Now().UTC()
	return &domain.Client{Email: b.em, PasswordHash: hp, FirstName: b.fn, LastName: b.ln, Company: b.co, Country: b.cnt, Currency: b.cur, Status: b.st, CreatedAt: n, UpdatedAt: n}, nil
}
