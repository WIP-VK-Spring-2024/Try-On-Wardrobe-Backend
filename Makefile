GO=go
BUILD_DIR=build
ENV=env POSTGRES.PASSWORD=12345

INTERNAL=internal/pkg
DOMAIN_PKG=${INTERNAL}/domain
DELIVERY_PKG=$$(${GO} list -f '{{.Dir}}' ./... | grep delivery | tr '\n' ' ')
ERRORS_PKG=${INTERNAL}/app_errors
PROTO_FILES=$$(find . -name *.proto)
GENERATED_DIR=internal/generated

.PHONY: easyjson sqlc gen build build_alpine run docker

easyjson:
		easyjson -snake_case -omit_empty -pkg ${DOMAIN_PKG} ${DELIVERY_PKG} ${ERRORS_PKG}

sqlc:
	sqlc generate

protoc:
	protoc --go-grpc_opt=paths=source_relative --go-grpc_out=${GENERATED_DIR} --go_opt=paths=source_relative --go_out=${GENERATED_DIR} ${PROTO_FILES}

gen: protoc sqlc easyjson 
	${GO} generate ./...

build:
	mkdir -p ${BUILD_DIR}
	${GO} build -o ${BUILD_DIR} ./...

build_alpine:
	mkdir -p ${BUILD_DIR}/alpine
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GO} build -o ${BUILD_DIR}/alpine ./...

run: build
	${ENV} ./${BUILD_DIR}/cmd

docker: build_alpine
	docker compose up --build -d
	docker compose logs app centrifugo --follow
