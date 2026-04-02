MODULE_NAME = omniport-api
BIN_DIR     = bin

.PHONY: dev tidy vet swagger pia build-monolith build-adm build-plan build-rls build-bill build-mrep build-all docker-monolith docker-adm

dev:
	go run ./cmd/monolith

build-monolith:
	go build -o $(BIN_DIR)/omniport-monolith ./cmd/monolith

build-adm:
	go build -o $(BIN_DIR)/omniport-adm ./cmd/adm-service

build-plan:
	go build -o $(BIN_DIR)/omniport-plan ./cmd/plan-service

build-rls:
	go build -o $(BIN_DIR)/omniport-rls ./cmd/rls-service

build-bill:
	go build -o $(BIN_DIR)/omniport-bill ./cmd/bill-service

build-mrep:
	go build -o $(BIN_DIR)/omniport-mrep ./cmd/mrep-service

build-all: build-monolith build-adm build-plan build-rls build-bill build-mrep

docker-monolith:
	docker build --build-arg ENTRY=monolith -t omniport:monolith .

docker-adm:
	docker build --build-arg ENTRY=adm-service -t omniport:adm .

tidy:
	go mod tidy

vet:
	go vet ./...

swagger:
	go run ./cmd/docs/swagger

pia:
	go run ./cmd/docs/pia
