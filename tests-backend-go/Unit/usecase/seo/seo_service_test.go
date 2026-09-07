package seo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/seo"
)

func TestSEOService_SitemapAndInfo(t *testing.T) {
	ctx := context.Background()
	pageRepo := memory.NewMockPageRepository()
	newsRepo := memory.NewMockNewsRepository()
	prodRepo := memory.NewMockProductRepository()

	// Seed a page
	_ = pageRepo.Create(ctx, &domain.Page{
		Title:     "About Us",
		Slug:      "about-us",
		Content:   "About FOSSBilling",
		Published: true,
	})

	svc := seo.NewSEOService(pageRepo, newsRepo, prodRepo)

	// 1. Get SEO Info
	info := svc.GetInfo("https://demo.fossbilling.org")
	if info.SitemapURL != "https://demo.fossbilling.org/sitemap.xml" {
		t.Errorf("Expected sitemap URL 'https://demo.fossbilling.org/sitemap.xml', got '%s'", info.SitemapURL)
	}

	// 2. Generate XML Sitemap
	xml, err := svc.GenerateSitemapXML(ctx, "https://demo.fossbilling.org")
	if err != nil {
		t.Fatalf("Failed to generate XML sitemap: %v", err)
	}

	if !strings.Contains(xml, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Errorf("Missing XML declaration in sitemap")
	}
	if !strings.Contains(xml, `<loc>https://demo.fossbilling.org/</loc>`) {
		t.Errorf("Missing home page in sitemap")
	}
	if !strings.Contains(xml, `about-us</loc>`) {
		t.Errorf("Missing published custom page in sitemap")
	}

	// 3. Ping Search Engines
	res, err := svc.PingSearchEngines(ctx, "https://demo.fossbilling.org")
	if err != nil {
		t.Fatalf("Failed to ping search engines: %v", err)
	}
	if len(res) == 0 {
		t.Errorf("Expected ping results, got empty map")
	}
}
