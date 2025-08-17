DB_PG_SCHEMA=amartha_loan_service

protocgen:
	cd proto/ && \
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false *.proto -I${GOPATH}/src -I. && \
	cd ..

# requires bufbuild/buf/buf
buf-gen:
	cd proto/ && \
	buf generate && \
	cd ..

mockgen:
	mockery --config mockery.yaml

migrate-up:
	migrate -path internal/migration -database "postgres://${DB_PG_HOST}:${DB_PG_PORT}/${DB_PG_SCHEMA}?sslmode=disable" up

migrate-down:
	migrate -path internal/migration -database "postgres://${DB_PG_HOST}:${DB_PG_PORT}/${DB_PG_SCHEMA}?sslmode=disable" down

docker-build:
	docker build -t authentication . && \
	docker tag authentication $(USERNAME)/authentication:$(VERSION) && \
	docker push $(USERNAME)/authentication:$(VERSION)

test:
	go clean -testcache && \
	go test -cover -race ./...

.PHONY: gen test protocgen