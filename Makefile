postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=ady -e POSTGRES_PASSWORD=password -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=ady --owner=ady simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose down


.PHONY: postgres createdb dropdb