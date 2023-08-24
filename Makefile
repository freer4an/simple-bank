postgres:
	docker run --name postgres-docker -p 5432:5432 --network bank-network -e POSTGRES_USER=ady -e POSTGRES_PASSWORD=password -d postgres:alpine

createdb:
	docker exec -it postgres-docker createdb --username=ady --owner=ady simple_bank

dropdb:
	docker exec -it postgres-docker dropdb simple_bank

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

dockerbuild:
	docker build -t simple-bank-image .

dockerrun:
	docker run --name simple-bank --network bank-network -p 8000:8000 -e DB_SOURCE=postgresql://ady:password@postgres-docker:5432/simple_bank?sslmode=disable -e REDIS_SOURCE=redis-docker:6379 simple-bank-image

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/freer4an/simple-bank/db/sqlc Store
	
evans:
	evans --host localhost --port 9000 -r repl

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

redis:
	docker run --name redis-docker -p 6379:6379 -d redis:alpine
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test dockerbuild dockerrun mock proto redis