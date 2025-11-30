# Go Production-Ready Microservice

> **🚀 Skip the setup, focus on logic: A production-ready Go microservice template with full observability built-in. ✨**

This project provides a robust, high-performance foundation for building scalable backend systems in Go (Golang). It adheres to **12-Factor App** principles and comes pre-configured with a complete monitoring stack, so you can start writing business logic immediately.

## ⚡ Key Features

* **Production-Grade Architecture:** Built on **Fiber** for speed, with clean code principles and graceful shutdown.
* **Full Observability:** Pre-configured **Prometheus** metrics and a professional **Grafana** Command Center dashboard.
* **DevOps Ready:** Ultra-lightweight Docker images (~15MB), multi-stage builds, and orchestrated via `docker-compose`.
* **Dynamic Configuration:** Environment-aware settings (`dev`, `prod`) with hot-reloading capabilities.
* **Secure by Default:** Runs as a non-root user inside containers.

## 🛠️ Quick Start

Get the entire stack (App + Prometheus + Grafana) running with a single command:

```bash
docker-compose up -d --build
```

### Access Points

* **API:** [http://localhost:3000](http://localhost:3000)
* **Grafana:** [http://localhost:3001](http://localhost:3001) (User: `admin` / Pass: `admin`)
* **Prometheus:** [http://localhost:9090](http://localhost:9090)

## 📡 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check (Returns 200 OK) |
| `GET` | `/error` | **Test Route:** Triggers a 500 error to test monitoring alerts |
| `GET` | `/metrics` | Exposes Prometheus metrics |

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details. You are free to use, modify, and distribute this software as you wish.
