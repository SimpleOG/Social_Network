db_version=postgres:latest
postgres_name=postgres
migrate_dir_path=..\..\..\internal\repositories\database\postgresql\migrations
postgres_user=root
postgres_password=secret
db_name=database
db_source=postgresql://$(postgres_user):$(postgres_password)@0.0.0.0:5050/$(db_name)?sslmode=disable
setup_docker_bd: pull_image run_postgres

pull_image:
	docker pull $(db_version)
run_postgres:
	docker run --name $(postgres_name)  -e POSTGRES_USER=$(postgres_user) -e POSTGRES_PASSWORD=$(postgres_password) -p 5050:5050 -d $(db_version)

createdb:
	docker exec -it $(postgres_name) createdb --username=$(postgres_user)  --owner=$(postgres_user)  $(db_name)
dropdb:
	docker exec -it $(postgres_name) dropdb $(db_name)

mg_create:
	migrate create -ext sql -dir $(migrate_dir_path) -seq init_schema
mg_up:
	migrate -path $(db_source) -verbose up
mg_down:
	migrate -path $(db_source) -verbose down

sqlc:
	cd ../sqlc && sqlc generate

#sqlc_install:
 	#go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest