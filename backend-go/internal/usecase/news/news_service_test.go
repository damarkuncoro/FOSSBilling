package news_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/news"
)

func TestNewsService_Lifecycle(t *testing.T) {
	repo := memory.NewMockNewsRepository()
	service := news.NewNewsService(repo)
	ctx := context.Background()

	// 1. Create Draft Post
	post, err := service.Create(ctx, news.CreateNewsDTO{
		AdminID: 1,
		Title:   "Scheduled Network Maintenance in Singapore DC",
		Content: "We will perform core router upgrades on Sunday 02:00 AM UTC.",
		Status:  domain.NewsStatusDraft,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if post.Slug != "scheduled-network-maintenance-in-singapore-dc" {
		t.Errorf("unexpected slug: %s", post.Slug)
	}

	// 2. Draft should not appear in public published list
	published, total, _ := service.ListPublished(ctx, 10, 0)
	if len(published) != 0 || total != 0 {
		t.Errorf("expected 0 published posts, got %d", total)
	}

	// 3. Publish Post
	updated, err := service.Update(ctx, post.ID, news.UpdateNewsDTO{
		Status: domain.NewsStatusPublished,
	})
	if err != nil || updated.Status != domain.NewsStatusPublished {
		t.Fatalf("Publish failed: %v", err)
	}

	// 4. Now visible in published list & fetch by slug
	published, total, _ = service.ListPublished(ctx, 10, 0)
	if len(published) != 1 || total != 1 {
		t.Errorf("expected 1 published post, got %d", total)
	}

	bySlug, err := service.GetBySlug(ctx, "scheduled-network-maintenance-in-singapore-dc")
	if err != nil || bySlug.ID != post.ID {
		t.Fatalf("GetBySlug failed: %v", err)
	}

	// 5. Delete Post
	err = service.Delete(ctx, post.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
