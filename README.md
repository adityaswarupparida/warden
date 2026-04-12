# Warden — Service Health Monitoring API

A service health monitoring application built with Go, PostgreSQL, and Kubernetes.

## Architecture
<img width="1567" height="785" alt="Screenshot 2026-04-12 at 1 29 58 PM" src="https://github.com/user-attachments/assets/1c0975bc-7c6c-4720-b47a-cc580270c621" />


## Architecture Choices
- Chose `PostgreSQL` for its *strong ACID guarantees, schema-on-write validation,* and *durable source of truth*. Native enum types enforce valid status values at the database level, rejecting incorrect writes before they reach the application. Most importantly, transactions make `incident create/resolve` and the corresponding `service status update` atomic guaranteeing the two are never out of sync.
- Opted for `pgxpool` over `pgBouncer` as it offers a *concurrency-safe, Go-native approach, maintaining a bounded pool of persistent connections per pod, allowing each pod to handle bursts independently without contention.* At current scale of up to 10-15 replicas, this handles connection management effectively, well within PostgreSQL's hard limits. pgxpool also *caches prepared statements per connection*, avoiding repeated query parsing and giving a performance edge over pgBouncer which breaks prepared statements in transaction mode.
- Chose `Sealed Secrets` over environment variables or plain `Kubernetes Secrets` for credential management, encrypting the PostgreSQL URL *asymmetrically* so it can be safely committed to git. Only the Sealed Secrets controller holding the *private key* can decrypt it, making the sealed secret useless to anyone without cluster access. This keeps secrets version-controlled alongside the rest of the infrastructure without exposing sensitive values..
- Went with `Prometheus` and `Grafana` as they are the *de facto standard for Kubernetes observability* — free, self-hosted, and battle-tested. Prometheus integrates natively with Kubernetes through ServiceMonitor CRDs, using service discovery to scrape `/metrics` from pods *dynamically*, handling the ephemeral nature of pods without any manual target configuration. PromQL makes it straightforward to express the four Golden Signals — traffic, errors, latency, and saturation. Grafana sits on top with a rich ecosystem of community dashboards and deep Prometheus integration. 
- Favoured `Kind` (Kubernetes in Docker) for local setup over managed services like `EKS` or `GKE` as it *provisions a full multi-node cluster entirely within Docker* — no cloud account, no cost, minimal overhead. Despite running locally, it *mirrors a production Kubernetes environment* closely enough to validate manifests, Helm charts, and service discovery before any real deployment. 


## Data Modelling Rationale and Connection Management Strategy

Designed a minimal schema around services and incident management.                                                                                                     
- **services** — the entity being monitored. Has a `status` enum (`healthy`, `degraded`, `down`) that reflects current health derived from active incidents, never set manually.
- **incidents** — events reported against a service. Has a `severity` enum (`low`, `medium`, `high`, `critical`), lifecycle timestamps (`created_at`, `updated_at`, `resolved_at`, `occurred_at`), and a `status` enum for incident lifecycle state. `resolved_at` is the source of truth for active incidents. Service status is updated atomically in the same transaction to *guarantee consistency*.
- **incident_updates** — audit trail that tracks the history of changes made to an incident (status changes, messages) as an incident progresses through its lifecycle (schema in place, not yet exposed via API).


Used **pgxpool** (application-level connection pooling) which maintains a pool of persistent connections:
- Connections are reused across requests — no TCP handshake overhead per request
- Pool size is bounded between `MinConns = 2` (always warm) and `MaxConns = 5` (hard cap)
- Though stateful, per-pod pooling gives better performance under load — up to 10-15 replicas (50-75 total connections) each pod handles bursts independently without contention.
- If all connections are busy, requests wait rather than failing immediately
- `pool.Ping()` on startup validates the connection before serving traffic
- Each connection in the pool maintains its own prepared statement cache — with MaxConns = 5, the same query can be parsed up to 5 times. Under normal load this is negligible, but worth noting for high-throughput scenarios with complex queries.


*Note: Beyond 15-20 replicas, pgxpool's linear connection growth becomes the bottleneck. In contrast, pgBouncer, which sits between the application and PostgreSQL, multiplexes hundreds of application connections into a small, fixed pool of database connections. As a result, PostgreSQL only sees what pgBouncer allows through regardless of how many pods are running. The tradeoff: it adds operational complexity and breaks prepared statements in transaction mode.*


## Prerequisites

Install the following tools before getting started:

- [Docker](https://docs.docker.com/get-docker/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) — `go install sigs.k8s.io/kind@latest`
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm](https://helm.sh/docs/intro/install/)
- [kubeseal](https://github.com/bitnami-labs/sealed-secrets#installation) — only required if bringing your own PostgreSQL URL (Option B)
- [cloud-provider-kind](https://github.com/kubernetes-sigs/cloud-provider-kind) — `go install sigs.k8s.io/cloud-provider-kind@latest` (for LoadBalancer support)

## Local Deployment

### 1. Clone the repository

```bash
git clone https://github.com/adityaswarupparida/warden
cd service-monitoring
```

### 2. Start cloud-provider-kind (in a separate terminal)

```bash
make cloud-provider
```

Keep this running in the background — it provides LoadBalancer IP assignment for Kind clusters.

### 4. Deploy the stack

Choose one of the two options below depending on your setup.

#### Option A — Using the provided sealed secrets key (in-cluster PostgreSQL)

Move the sealed secrets key into the project root:

```bash
mv /path/to/sealed-secrets-key.yml .
```

> Never commit `sealed-secrets-key.yml` to git.

Then run:

```bash
make setup
```

This will:
- Create a Kind cluster
- Install Sealed Secrets and restore the encryption key
- Install PostgreSQL
- Install Prometheus + Grafana monitoring stack
- Deploy the warden application

#### Option B — Using your own PostgreSQL URL (external DB, managed service, Docker, etc.)

Set your PostgreSQL URL in `.env`:

```
POSTGRESQL_URL=postgres://<user>:<password>@<host>:5432/<db>?sslmode=disable
```

Then run:

```bash
make setup-external
```

This will:
- Create a Kind cluster
- Install Sealed Secrets and seal your URL from `.env`
- Install Prometheus + Grafana monitoring stack
- Deploy the warden application

### 6. Get the service IP

```bash
kubectl get svc warden
```

Use the `EXTERNAL-IP` as your base URL:

```bash
export BASE=http://<EXTERNAL-IP>
```

---

## API Usage

### Services

**Create a service**

```bash
curl -X POST $BASE/api/v1/services \
  -H "Content-Type: application/json" \
  -d '{"name":"payments","description":"payment processing service"}'
```

**List all services**

```bash
curl $BASE/api/v1/services
```

---

### Incidents

**Create an incident**

```bash
curl -X POST $BASE/api/v1/services/<service_id>/incidents \
  -H "Content-Type: application/json" \
  -d '{"title":"DB latency spike","description":"p99 latency above 2s","severity":"high","occurred_at":"2026-04-12T10:00:00Z"}'
```

Severity options: `low`, `medium`, `high`, `critical`

> `occurred_at` is optional — use it to record when the issue actually started, which may differ from when it was reported.

**Get incidents for a service**

```bash
curl $BASE/api/v1/services/<service_id>/incidents
```

**Resolve an incident**

```bash
curl -X PATCH $BASE/api/v1/services/<service_id>/incidents/<incident_id>/resolve
```

---

## Load Testing

Generate traffic to populate Grafana dashboards:

```bash
# GET requests
for i in {1..500}; do curl -s $BASE/api/v1/services > /dev/null; sleep 0.2; done

# POST requests
for i in {1..50}; do
  curl -s -X POST $BASE/api/v1/services \
    -H "Content-Type: application/json" \
    -d '{"name":"service-'$i'","description":"test service '$i'"}' > /dev/null
  sleep 0.2
done

# Generate errors with invalid requests (blank name and description)
for i in {1..50}; do
  curl -s -X POST $BASE/api/v1/services \
    -H "Content-Type: application/json" \
    -d '{"name":"","description":""}' > /dev/null
  sleep 0.2
done
```

---

## Health & Metrics

```bash
curl $BASE/healthz
curl $BASE/metrics
```

---

## Observability

**Access Grafana**

```bash
kubectl get svc monitoring-grafana
```

Open `http://<EXTERNAL-IP>` in your browser. Login with `admin` / `admin`.

The **Service Health Monitoring** dashboard visualises the four Golden Signals:
- Traffic — requests per second
- Errors — error rate
- Latency — p50, p95, p99
- Saturation — CPU and memory usage

---

## Infrastructure as Code

[Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) is used as the IaC tool to provision the local Kubernetes cluster. The cluster configuration is defined in `deploy/cluster.yml` (1 control plane + 2 worker nodes) and the entire stack is orchestrated via `Makefile` targets — making the setup fully reproducible and automated.

---

## Delete Cluster

```bash
make delete
```
