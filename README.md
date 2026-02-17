# Go Production-Ready Microservice

> **🚀 Skip the setup, focus on logic: A production-ready Go microservice template with full observability built-in. ✨**

This project provides a robust, high-performance foundation for building scalable backend systems in Go (Golang). It adheres to **12-Factor App** principles and comes pre-configured with a complete monitoring stack, so you can start writing business logic immediately.

## ⚡ Key Features

* **Production-Grade Architecture:** Built on **Fiber** for speed, with clean code principles and graceful shutdown.
* **Full Observability:** Pre-configured **Prometheus** metrics and a professional **Grafana** Command Center dashboard.
* **Production-Grade Architecture:** Built on **Fiber** for speed, with clean code principles and graceful shutdown.
* **Full Observability:** Pre-configured **Prometheus** metrics and a professional **Grafana** Command Center dashboard.
* **DevOps Ready:** Ultra-lightweight Docker images (~15MB), multi-stage builds, and orchestrated via `docker-compose`.
* **Dynamic Configuration:** Environment-aware settings (`dev`, `prod`) with hot-reloading capabilities.
* **Secure by Default:** Runs as a non-root user inside containers.

## 🛠️ Quick Start

**1. Start Everything (App + Monitoring):**

```bash
make all
```

**2. Start Only the App:**

```bash
make up
```

**3. Start Only Monitoring:**

```bash
make infra
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

## 👩‍💻 Development Guide

### Adding a New Endpoint

Since the project uses **Fiber** and follows a **Hexagonal Architecture**, here is the standard flow to add a new feature:

1. **Define the Interface (Port):**
    * Create a new method in the `Service` interface (`internal/service/service.go`) if needed.
    * Create a new method in the `Repository` interface (`internal/repository/repository.go`) if database access is required.

2. **Implement the Business Logic (Service):**
    * Implement the method in `internal/service/user_service.go` (or create a new service).

3. **Implement the Handler (Adapter):**
    * Create a new handler method in `internal/transport/http/handler/`.
    * Use `fiber.Ctx` to get inputs and return responses.

4. **Register the Route:**
    * Add the new route in `internal/transport/http/router/router.go`.
    * Bind it to the handler method.

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details. You are free to use, modify, and distribute this software as you wish.
