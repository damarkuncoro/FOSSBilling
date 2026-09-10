package payment

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
)

type PaymentService struct {
	reg *payment.GatewayRegistry; ir domain.InvoiceRepository; cr domain.ClientRepository
}

func NewPaymentService(r *payment.GatewayRegistry, ir domain.InvoiceRepository, cr domain.ClientRepository) *PaymentService {
	return &PaymentService{r, ir, cr}
}

func (s *PaymentService) InitiateInvoicePayment(ctx context.Context, id int64, gid string) (*payment.PaymentResponse, error) {
	inv, err := s.ir.GetByID(ctx, id); if err != nil { return nil, err }
	c, err := s.cr.GetByID(ctx, inv.ClientID); if err != nil { return nil, err }
	gw, err := s.reg.Get(gid); if err != nil { return nil, err }

	return gw.InitiatePayment(ctx, payment.PaymentRequest{
		InvoiceID: inv.ID, InvoiceNr: inv.Serie + inv.Nr, Amount: inv.Total, Currency: inv.Currency,
		ClientEmail: c.Email, ClientName: c.FirstName + " " + c.LastName,
		ReturnURL: "/client/invoice/" + fmt.Sprint(inv.ID), CancelURL: "/client/invoice/" + fmt.Sprint(inv.ID),
	})
}

func (s *PaymentService) ListAvailableGateways() []payment.PaymentGateway { return s.reg.List() }
