db_version=postgres:latest
postgres_name=postgres
migrate_dir_path=../../../internal/repositories/database/postgresql/migrations/
postgres_user=root
postgres_password=1234
db_name=database
db_port=5050

db_source=postgresql://$(postgres_user):$(postgres_password)@0.0.0.0:$(db_port)/$(db_name)?sslmode=disable

setup_docker_bd: pull_image run_postgres

pull_image:
	docker pull $(db_version)
create_postgres_container:
	docker run --name $(postgres_name)  -e POSTGRES_USER=$(postgres_user) -e POSTGRES_PASSWORD=$(postgres_password) -p $(db_port):5432 -d $(db_version)

createdb:
	docker exec -it $(postgres_name) createdb --username=$(postgres_user)  --owner=$(postgres_user)  $(db_name)
dropdb:
	docker exec -it $(postgres_name) dropdb $(db_name)

mg_create:
	migrate create -ext sql -dir $(migrate_dir_path) -seq init_schema_1
mg_up:
	migrate -path $(migrate_dir_path) -database $(db_source) -verbose up
mg_down:
	migrate -path $(migrate_dir_path) -database $(db_source)  -verbose down

sqlc:
	cd ../sqlc && sqlc generate

#sqlc_install:
 	#go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

pull_redis:
	docker pull redis:latest
create_redis:
	docker run --name redis -p 6379:6379 -d redis
