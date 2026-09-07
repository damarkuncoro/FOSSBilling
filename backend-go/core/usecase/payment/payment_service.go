package payment

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
)

type PaymentService struct {
	gatewayRegistry *payment.GatewayRegistry
	invoiceRepo     domain.InvoiceRepository
	clientRepo      domain.ClientRepository
}

func NewPaymentService(
	gatewayRegistry *payment.GatewayRegistry,
	invoiceRepo domain.InvoiceRepository,
	clientRepo domain.ClientRepository,
) *PaymentService {
	return &PaymentService{
		gatewayRegistry: gatewayRegistry,
		invoiceRepo:     invoiceRepo,
		clientRepo:      clientRepo,
	}
}

func (s *PaymentService) InitiateInvoicePayment(ctx context.Context, invoiceID int64, gatewayID string) (*payment.PaymentResponse, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepo.GetByID(ctx, inv.ClientID)
	if err != nil {
		return nil, err
	}

	gw, err := s.gatewayRegistry.Get(gatewayID)
	if err != nil {
		return nil, err
	}

	req := payment.PaymentRequest{
		InvoiceID:   inv.ID,
		InvoiceNr:   fmt.Sprintf("%s%s", inv.Serie, inv.Nr),
		Amount:      inv.Total,
		Currency:    inv.Currency,
		ClientEmail: client.Email,
		ClientName:  fmt.Sprintf("%s %s", client.FirstName, client.LastName),
		ReturnURL:   "/client/invoice/" + fmt.Sprint(inv.ID),
		CancelURL:   "/client/invoice/" + fmt.Sprint(inv.ID),
	}

	return gw.InitiatePayment(ctx, req)
}

func (s *PaymentService) ListAvailableGateways() []payment.PaymentGateway {
	return s.gatewayRegistry.List()
}
