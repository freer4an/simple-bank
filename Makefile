postgres:
	docker run --name postgres -p 5432:5432 --network bank-network -e POSTGRES_USER=ady -e POSTGRES_PASSWORD=password -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=ady --owner=ady simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
	
server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/freer4an/simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup, migratedown, sqlc, test, server, mock