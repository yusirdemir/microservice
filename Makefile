.PHONY: up down infra infra-down all clean help

help:
	@echo "Available commands:"
	@echo "  make up          - Start the application (User Service)"
	@echo "  make down        - Stop the application"
	@echo "  make infra       - Start the monitoring stack (Prometheus + Grafana)"
	@echo "  make infra-down  - Stop the monitoring stack"
	@echo "  make all         - Start everything (App + Infra)"
	@echo "  make clean       - Stop everything"

up:
	docker-compose up -d --build

down:
	docker-compose down

infra:
	docker-compose -f deploy/docker-compose.infra.yml up -d

infra-down:
	docker-compose -f deploy/docker-compose.infra.yml down

all: up infra

clean: down infra-down
