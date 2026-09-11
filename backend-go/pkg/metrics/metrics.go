package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "fossbilling_http_requests_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "path", "status"})

	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fossbilling_http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	InvoicesPaidTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fossbilling_invoices_paid_total",
		Help: "Total number of invoices marked as paid",
	})

	RevenueTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "fossbilling_revenue_total",
		Help: "Total revenue collected in minor units",
	}, []string{"currency"})

	ActiveOrders = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fossbilling_active_orders_count",
		Help: "Current count of active service orders",
	})

	ProvisioningFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fossbilling_provisioning_failures_total",
		Help: "Total number of failed provisioning attempts",
	})
)
