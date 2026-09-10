# FOSSBilling Kubernetes Deployment

This directory contains the base manifests for deploying FOSSBilling on a Kubernetes cluster.

## Architecture

* **API Deployment:** Scalable Go REST API with health probes.
* **Worker Deployment:** Background task processor for cron and events.
* **PostgreSQL:** Managed via StatefulSet or persistent Volume (Base example uses a simple Deployment).
* **Redis:** Used for caching and event bus synchronization.

## Getting Started

1. **Configure Secrets & ConfigMaps:**
   Copy `base/configmap.yaml.example` and `base/secrets.yaml.example` (not provided, create based on `.env`).

2. **Apply Manifests:**
   ```bash
   kubectl apply -f base/
   ```

3. **Expose Services:**
   Configure an Ingress controller (Nginx/Traefik) using the provided `base/ingress.yaml`.

## Production Recommendations

* Use a managed database service (Amazon RDS, Google Cloud SQL) instead of in-cluster Postgres.
* Use Horizontal Pod Autoscaler (HPA) for the API deployment.
* Configure TLS/SSL using Cert-manager.
