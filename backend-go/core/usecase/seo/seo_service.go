package seo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type SEOService struct {
	pageRepo domain.PageRepository; newsRepo domain.NewsRepository; productRepo domain.ProductRepository; httpClient *http.Client; lastPingAt *time.Time
}

func NewSEOService(pr domain.PageRepository, nr domain.NewsRepository, pdr domain.ProductRepository) *SEOService {
	return &SEOService{pr, nr, pdr, &http.Client{Timeout: 10 * time.Second}, nil}
}

type SEOInfo struct { SitemapURL string `json:"sitemap_url"`; LastPingAt *time.Time `json:"last_ping_at"`; Engines map[string]string `json:"engines"` }

func (s *SEOService) GetInfo(bURL string) SEOInfo {
	bURL = strings.TrimSuffix(bURL, "/")
	return SEOInfo{fmt.Sprintf("%s/sitemap.xml", bURL), s.lastPingAt, map[string]string{"google": "https://www.google.com/ping?sitemap=", "bing": "https://www.bing.com/ping?sitemap="}}
}

func (s *SEOService) PingSearchEngines(ctx context.Context, bURL string) (map[string]bool, error) {
	info, res := s.GetInfo(bURL), make(map[string]bool)
	enc := url.QueryEscape(info.SitemapURL)
	for n, e := range info.Engines {
		req, _ := http.NewRequestWithContext(ctx, "GET", e+enc, nil)
		resp, err := s.httpClient.Do(req)
		res[n] = err == nil && (resp == nil || resp.StatusCode < 400)
		if resp != nil { resp.Body.Close() }
	}
	now := time.Now().UTC(); s.lastPingAt = &now
	return res, nil
}

func (s *SEOService) GenerateSitemapXML(ctx context.Context, bURL string) (string, error) {
	bURL = strings.TrimSuffix(bURL, "/")
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	add := func(loc, freq, pri, date string) { b.WriteString(fmt.Sprintf("  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n  </url>\n", loc, date, freq, pri)) }

	add(bURL+"/", "daily", "1.0", time.Now().Format("2006-01-02"))
	if ps, _, err := s.pageRepo.List(ctx, 100, 0); err == nil {
		for _, p := range ps { if p.Published { add(fmt.Sprintf("%s/pages/%s", bURL, p.Slug), "weekly", "0.8", p.UpdatedAt.Format("2006-01-02")) } }
	}
	if ns, _, err := s.newsRepo.ListPublished(ctx, 100, 0); err == nil {
		for _, n := range ns { add(fmt.Sprintf("%s/news/%s", bURL, n.Slug), "monthly", "0.7", n.UpdatedAt.Format("2006-01-02")) }
	}
	if pds, _, err := s.productRepo.List(ctx, 100, 0); err == nil {
		for _, p := range pds { if p.Status == "enabled" { add(fmt.Sprintf("%s/order?product=%d", bURL, p.ID), "weekly", "0.9", p.UpdatedAt.Format("2006-01-02")) } }
	}
	b.WriteString(`</urlset>`)
	return b.String(), nil
}
