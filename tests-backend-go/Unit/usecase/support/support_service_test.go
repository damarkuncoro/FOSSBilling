package support_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
)

func setupSupportService() (*support.SupportService, *memory.MockSupportRepository, *memory.MockClientRepository) {
	supportRepo := memory.NewMockSupportRepository()
	clientRepo := memory.NewMockClientRepository()
	service := support.NewSupportService(supportRepo, clientRepo)

	return service, supportRepo, clientRepo
}

func TestSupportService_OpenTicket(t *testing.T) {
	ctx := context.Background()
	service, _, clientRepo := setupSupportService()

	client := &domain.Client{
		Email:     "help.seeker@example.com",
		FirstName: "David",
	}
	_ = clientRepo.Create(ctx, client)

	dto := support.CreateTicketDTO{
		ClientID:   client.ID,
		HelpdeskID: 1,
		Subject:    "Cannot connect to MySQL database",
		Message:    "I keep getting connection timeout when connecting to port 3306.",
		Priority:   domain.PriorityHigh,
	}

	ticket, err := service.OpenTicket(ctx, dto)
	if err != nil {
		t.Fatalf("OpenTicket failed: %v", err)
	}

	if ticket.Status != domain.TicketStatusOpen {
		t.Errorf("Status = %s; want open", ticket.Status)
	}
	if ticket.Priority != domain.PriorityHigh {
		t.Errorf("Priority = %s; want high", ticket.Priority)
	}
}

func TestSupportService_ConversationThreadingAndStatusTransitions(t *testing.T) {
	ctx := context.Background()
	service, supportRepo, clientRepo := setupSupportService()

	client := &domain.Client{
		Email:     "client@example.com",
		FirstName: "Sarah",
	}
	_ = clientRepo.Create(ctx, client)

	// 1. Open Ticket
	ticket, _ := service.OpenTicket(ctx, support.CreateTicketDTO{
		ClientID:   client.ID,
		HelpdeskID: 1,
		Subject:    "SSL Certificate Installation",
		Message:    "Please install Let's Encrypt SSL on my domain.",
	})

	// 2. Staff Replies -> Status changes to awaiting_client
	adminID := int64(1)
	staffMsg, err := service.StaffReply(ctx, ticket.ID, adminID, "SSL has been issued. Please clear your browser cache.")
	if err != nil {
		t.Fatalf("StaffReply failed: %v", err)
	}
	if staffMsg.AdminID == nil || *staffMsg.AdminID != adminID {
		t.Errorf("StaffMsg AdminID = %v; want %d", staffMsg.AdminID, adminID)
	}

	updatedTicket1, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	if updatedTicket1.Status != domain.TicketStatusAwaitingClient {
		t.Errorf("After staff reply, status = %s; want awaiting_client", updatedTicket1.Status)
	}

	// 3. Client Replies -> Status changes to awaiting_staff
	_, err = service.ClientReply(ctx, ticket.ID, client.ID, "Thank you! It works now.", "127.0.0.1")
	if err != nil {
		t.Fatalf("ClientReply failed: %v", err)
	}

	updatedTicket2, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	if updatedTicket2.Status != domain.TicketStatusAwaitingStaff {
		t.Errorf("After client reply, status = %s; want awaiting_staff", updatedTicket2.Status)
	}

	// 4. Close Ticket
	_ = service.CloseTicket(ctx, ticket.ID, 0)
	closedTicket, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	if closedTicket.Status != domain.TicketStatusClosed {
		t.Errorf("Closed status = %s; want closed", closedTicket.Status)
	}

	// 5. Reply to closed ticket should fail
	_, replyClosedErr := service.ClientReply(ctx, ticket.ID, client.ID, "Another question...", "127.0.0.1")
	if replyClosedErr != support.ErrTicketClosed {
		t.Errorf("Expected ErrTicketClosed, got: %v", replyClosedErr)
	}
}
