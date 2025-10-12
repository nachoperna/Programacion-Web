# Variables de configuración
DB_URL=postgres://nachoperna:nachobdtpe@localhost:5432/BD_TPEspecial?sslmode=disable
MIGRATIONS_DIR=./db/migrations
DB_CONTAINER=BD_TPEspecial
DB_USER=nachoperna
DB_NAME=BD_TPEspecial

# -----------------------------
# Comandos principales
# -----------------------------

# Generar código con sqlc
sqlc:
	sqlc generate

# Levantar servicios con Docker (en segundo plano)
up:
	docker-compose up -d

# Apagar servicios y borrar volúmenes (datos)
down:
	docker-compose down -v

# Esperar a que Postgres esté listo
wait:
	@echo "Esperando a que Postgres arranque..."
	@until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Postgres está listo ✅"

# Crear una nueva migración
# Uso: make migrate-create name=nombre_migracion
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Ejecutar todas las migraciones hacia arriba
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Revertir la última migración
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Revertir todas las migraciones
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -all

# Reconstruir la base desde cero: baja, sube, espera, aplica migraciones y genera código
reset: down up wait migrate-up sqlc

.PHONY: sqlc up down wait migrate-create migrate-up migrate-down migrate-reset reset

