# Custom Payment Gateways

FOSSBilling is designed to be extensible, allowing you to add support for any payment provider by implementing a simple interface in Go.

---

## 🛠️ The Gateway Interface

Every payment gateway must implement the `PaymentGateway` interface found in `core/domain/payment.go`:

```go
type PaymentGateway interface {
    ID() string
    Name() string
    InitiatePayment(ctx context.Context, invoice *Invoice) (*PaymentSession, error)
    HandleWebhook(ctx context.Context, payload []byte) (*TransactionResult, error)
}
```

### Methods
- **ID():** Returns a unique slug (e.g., `"stripe"`).
- **Name():** Returns a human-readable display name.
- **InitiatePayment():** Returns a redirect URL or instructions for the client to pay.
- **HandleWebhook():** Processes IPNs (Instant Payment Notifications) from the provider to confirm payment.

---

## 🚀 Creating Your Own Gateway

Follow these steps to add a new provider:

1. **Create Package:** Create a new file in `core/service/payment/gateways/yourprovider.go`.
2. **Implement Struct:** Define a struct to hold your credentials and configuration.
3. **Register Gateway:** Add your new gateway to the `GatewayRegistry` in `cmd/api/services.go`:

```go
gatewayRegistry.Register(gateways.NewYourProviderGateway(config))
```

---

## 🔔 Webhook Security

When implementing `HandleWebhook`, always verify the origin and signature of the incoming request:
- **IP Validation:** Check if the request comes from known provider IP ranges.
- **Signature Check:** Use HMAC or SHA-256 signatures provided in the headers.
- **Invoice Match:** Ensure the amount and currency in the webhook match the invoice in your database.

---

## 💳 Supported Gateway Types

The FOSSBilling frontend can render different UIs based on the gateway type:
- **Redirect:** Sends client to an external site (PayPal, Stripe Checkout).
- **Embedded:** Renders a form directly in the portal.
- **Instructional:** Shows bank transfer details or QRIS codes.
