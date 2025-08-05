ENV ?= dev

service-build-gpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.gpu.yml build
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.gpu.yml build
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-build-cpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.cpu.yml build
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.cpu.yml build
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-up-gpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.gpu.yml up -d
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.gpu.yml up -d
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-up-cpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.cpu.yml up -d
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.cpu.yml up -d
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-down-gpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.gpu.yml down
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.gpu.yml down
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-down-cpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.cpu.yml down
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.cpu.yml down
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-clean-gpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.gpu.yml down -v
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.gpu.yml down -v
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

service-clean-cpu:
ifeq ($(ENV), dev)
	docker compose -f docker-compose.override.yml -f docker-compose.cpu.yml down -v
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.prod.yml -f docker-compose.cpu.yml down -v
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

