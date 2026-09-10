package catalog

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

type ServerService struct {
	repo domain.CatalogRepository; factory *provisioning.ProvisionerFactory
}

func NewServerService(r domain.CatalogRepository, f *provisioning.ProvisionerFactory) *ServerService {
	return &ServerService{r, f}
}

func (s *ServerService) ListServers(ctx context.Context) ([]*domain.Server, error) { return s.repo.ListServers(ctx) }

func (s *ServerService) CreateServer(ctx context.Context, srv *domain.Server) error {
	if srv.Status == "" { srv.Status = "online" }; return s.repo.CreateServer(ctx, srv)
}

func (s *ServerService) DeleteServer(ctx context.Context, id int64) error { return s.repo.DeleteServer(ctx, id) }

func (s *ServerService) TestConnection(ctx context.Context, id int64) error {
	srv, err := s.repo.GetServerByID(ctx, id); if err != nil { return err }
	prov, err := s.factory.CreateProvisioner(provisioning.ServerConfig{Type: srv.Manager, Host: srv.Hostname, Username: "admin", APIToken: srv.AccessKey, UseSSL: true})
	if err != nil { return err }; return prov.TestConnection(ctx)
}
