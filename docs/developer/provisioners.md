# Custom Provisioners & Driver Development

FOSSBilling makes it straightforward to add custom provisioning drivers for cloud hypervisors, VPS panels, game servers, or SaaS APIs.

---

## 🛠️ The `ServiceProvisioner` Interface

All provisioning drivers implement the standard Go interface:

```go
type ServiceProvisioner interface {
    Type() domain.ProductType
    Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error)
    Suspend(ctx context.Context, order *domain.Order, reason string) error
    Unsuspend(ctx context.Context, order *domain.Order) error
    Terminate(ctx context.Context, order *domain.Order) error
    Renew(ctx context.Context, order *domain.Order) error
    Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error)
    ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error
}
```

---

## 📝 Registering a New Driver

1. Implement the interface in `backend-go/core/service/provisioning/<driver_name>.go`.
2. Register the driver in `ProvisionerRegistry` (`backend-go/core/service/provisioning/registry.go`).
3. Add driver unit tests in `tests-backend-go/Unit/service/provisioning/`.
