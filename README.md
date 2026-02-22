# Go Microservice Template

A production-ready, lightweight microservice template built with Go. This project is designed to be simple, fast, and fully observable.

## 🛠️ Stack & Technologies Used

- **Go (Golang):** The core programming language for high performance.
- **Fiber:** A fast web framework for Go, used for handling HTTP API requests.
- **Couchbase:** A NoSQL document database used for fast, scalable data storage.
- **Zap (Uber):** High-performance logging library used for structured, JSON-based application logs.
- **OpenTelemetry & Jaeger:** Used for distributed tracing (tracking a request as it moves through the system to find bottlenecks).
- **Prometheus:** Scrapes and stores time-series metrics (like CPU usage, memory, request counts).
- **Grafana:** Command center dashboard to visualize the data/metrics collected by Prometheus.
- **Docker & Docker Compose:** Used to containerize the app and easily start all infrastructure (database, monitoring) with a single command.

## 🚀 Quick Start

Start everything (Application + Database + Monitoring) with one command:

```bash
make all
```

### 🎯 Access Points

Once started, you can access the services here:

- **API:** [http://localhost:3000](http://localhost:3000)
- **Grafana (Metrics Dashboard):** [http://localhost:3001](http://localhost:3001) *(User: admin / Pass: admin)*
- **Jaeger (Tracing UI):** [http://localhost:16686](http://localhost:16686)
- **Prometheus (Raw Metrics):** [http://localhost:9090](http://localhost:9090)
- **Couchbase (Database UI):** [http://localhost:8091](http://localhost:8091)

## 🛑 Stop the Project

To stop and clean up everything:

```bash
make clean
```
