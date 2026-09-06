package gateways

import (
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
)

type GatewayConfig struct {
	ID       string                 `json:"id"`       // "midtrans", "stripe", "paypal", "bank_transfer"
	Enabled  bool                   `json:"enabled"`
	IsProd   bool                   `json:"is_prod"`
	Settings map[string]interface{} `json:"settings"`
}

// PaymentGatewayFactory constructs PaymentGateway driver instances from dynamic gateway configurations
type PaymentGatewayFactory struct{}

func NewPaymentGatewayFactory() *PaymentGatewayFactory {
	return &PaymentGatewayFactory{}
}

func (f *PaymentGatewayFactory) CreateGateway(cfg GatewayConfig) (payment.PaymentGateway, error) {
	getString := func(key string) string {
		if val, ok := cfg.Settings[key].(string); ok {
			return val
		}
		return ""
	}

	switch cfg.ID {
	case "midtrans":
		serverKey := getString("server_key")
		clientKey := getString("client_key")
		return NewMidtransGateway(serverKey, clientKey, cfg.IsProd), nil

	case "stripe":
		secretKey := getString("secret_key")
		pubKey := getString("publishable_key")
		whSec := getString("webhook_secret")
		return NewStripeGateway(secretKey, pubKey, whSec), nil

	case "paypal":
		clientID := getString("client_id")
		secret := getString("secret")
		return NewPayPalGateway(clientID, secret, cfg.IsProd), nil

	case "bank_transfer":
		bankName := getString("bank_name")
		accNr := getString("account_number")
		holder := getString("account_holder")
		return NewBankTransferGateway(bankName, accNr, holder), nil

	default:
		return nil, fmt.Errorf("unsupported payment gateway provider: %s", cfg.ID)
	}
}
