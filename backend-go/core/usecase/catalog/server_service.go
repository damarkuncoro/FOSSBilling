package catalog

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

type ServerService struct {
	catalogRepo domain.CatalogRepository
	factory     *provisioning.ProvisionerFactory
}

func NewServerService(catalogRepo domain.CatalogRepository, factory *provisioning.ProvisionerFactory) *ServerService {
	return &ServerService{
		catalogRepo: catalogRepo,
		factory:     factory,
	}
}

func (s *ServerService) ListServers(ctx context.Context) ([]*domain.Server, error) {
	return s.catalogRepo.ListServers(ctx)
}

func (s *ServerService) CreateServer(ctx context.Context, server *domain.Server) error {
	if server.Status == "" {
		server.Status = "online"
	}
	return s.catalogRepo.CreateServer(ctx, server)
}

func (s *ServerService) DeleteServer(ctx context.Context, id int64) error {
	return s.catalogRepo.DeleteServer(ctx, id)
}

func (s *ServerService) TestConnection(ctx context.Context, id int64) error {
	server, err := s.catalogRepo.GetServerByID(ctx, id)
	if err != nil {
		return err
	}

	// Instantiate temporary provisioner for testing
	cfg := provisioning.ServerConfig{
		Type:     server.Manager,
		Host:     server.Hostname,
		Username: "admin", // In production, these should be stored securely
		APIToken: server.AccessKey,
		UseSSL:   true,
	}

	prov, err := s.factory.CreateProvisioner(cfg)
	if err != nil {
		return err
	}

	return prov.TestConnection(ctx)
}
