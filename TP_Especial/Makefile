# Variables de configuración
DB_URL=postgres://nachoperna:nachobdtpe@localhost:5432/BD_TPEspecial?sslmode=disable
MIGRATIONS_DIR=./db/migrations
DB_CONTAINER=BD_TPEspecial
DB_USER=nachoperna
DB_NAME=BD_TPEspecial
MIGRATE_DOCKER = docker run --rm -v "$(PWD)/$(MIGRATIONS_DIR):/migrations:z" --network host migrate/migrate

# Instalación de dependencias 

install-docker:
	@echo "Instalando Docker 28.5.1..."
	curl -fsSL https://get.docker.com -o get-docker.sh
	sudo sh get-docker.sh
	sudo systemctl enable --now docker
	rm -f get-docker.sh
	@echo "Docker instalado."

install-go:
	@echo "Instalando Go 1.24.9..."
	curl -LO https://go.dev/dl/go1.24.9.linux-amd64.tar.gz
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xzf go1.24.9.linux-amd64.tar.gz
	rm go1.24.9.linux-amd64.tar.gz
	@echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.bashrc
	@echo "Go instalado."

install-templ:
	@echo "Instalando templ v0.3.960..."
	curl -LO https://github.com/a-h/templ/releases/download/v0.3.960/templ_Linux_x86_64.tar.gz
	tar -xzf templ_Linux_x86_64.tar.gz
	sudo mv templ /usr/local/bin/
	rm templ_Linux_x86_64.tar.gz
	@echo "Templ instalado."

install-sqlc:
	@echo "Instalando sqlc v1.30.0..."
	curl -LO https://github.com/sqlc-dev/sqlc/releases/download/v1.30.0/sqlc_1.30.0_linux_amd64.tar.gz
	tar -xzf sqlc_1.30.0_linux_amd64.tar.gz
	sudo mv sqlc /usr/local/bin/
	rm sqlc_1.30.0_linux_amd64.tar.gz
	@echo "sqlc instalado."

install-hurl:
	@echo "Instalando Hurl 4.2.0..."
	curl -LO https://github.com/Orange-OpenSource/hurl/releases/download/4.2.0/hurl-4.2.0-x86_64-linux.tar.gz
	tar -xzf hurl-4.2.0-x86_64-linux.tar.gz
	sudo mv hurl/bin/hurl /usr/local/bin/
	rm -rf hurl hurl-4.2.0-x86_64-linux.tar.gz
	@echo "Hurl instalado."


install-all: install-docker install-go install-templ install-sqlc install-hurl


# Comandos principales

# Generar código con sqlc
sqlc:
	sqlc generate

# Levantar servicios con Docker (en segundo plano)
up:
	docker-compose up -d

# Apagar servicios y borrar volúmenes (datos)
down:
	docker-compose down -v
	@PID=$$(lsof -t -i :8080); \
	if [ -n "$$PID" ]; then \
		echo "Proceso encontrado (PID: $$PID). Deteniendo..."; \
		kill $$PID; \
	else \
		echo "No se encontró ningún proceso en el puerto 8080."; \
	fi
	@echo "Servidor y puerto 8080 cerrado"

go_down:
	@PID=$$(lsof -t -i :8080); \
	if [ -n "$$PID" ]; then \
		echo "Proceso encontrado (PID: $$PID). Deteniendo..."; \
		kill $$PID; \
	else \
		echo "No se encontró ningún proceso en el puerto 8080."; \
	fi
	@echo "Servidor y puerto 8080 cerrado"

# Esperar a que Postgres esté listo
wait:
	@echo "Esperando a que Postgres arranque..."
	@until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Postgres está listo ✅"

# Crear una nueva migración
# Uso: make migrate-create name=nombre_migracion
migrate-up:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" down

migrate-create:
	$(MIGRATE_DOCKER) create -ext sql -dir /migrations -seq $(name)

# Revertir todas las migraciones
migrate-reset:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" down -all

templ:
	templ generate

inserts:
	hurl --test inserts.hurl

test:
	hurl --test tests.hurl

go:
	go run . &
	@sleep 2
# Reconstruir la base desde cero: baja, sube, espera, aplica migraciones y genera código
reset: down up wait migrate-up sqlc templ

run: reset go test 

.PHONY: sqlc up down wait migrate-create migrate-up migrate-down migrate-reset templ reset go run test install-docker install-go install-templ install-sqlc install-hurl install-all

