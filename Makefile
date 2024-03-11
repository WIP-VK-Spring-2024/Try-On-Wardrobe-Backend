GO=go
BUILD_DIR=build
ENV=env POSTGRES.PASSWORD=12345

INTERNAL=internal/pkg
DOMAIN_PKG=${INTERNAL}/domain
DELIVERY_PKG=$$(${GO} list -f '{{.Dir}}' ./... | grep delivery | tr '\n' ' ')
ERRORS_PKG=${INTERNAL}/app_errors

gen:
	easyjson -snake_case -omit_empty -pkg ${DOMAIN_PKG} ${DELIVERY_PKG} ${ERRORS_PKG}
	${GO} generate ./...

build: gen
	mkdir -p ${BUILD_DIR}
	${GO} build -o ${BUILD_DIR} ./...

build_alpine: gen
	mkdir -p ${BUILD_DIR}/alpine
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GO} build -o ${BUILD_DIR}/alpine ./...

run: build
	${ENV} ./${BUILD_DIR}/cmd

docker: build_alpine
	docker compose up --build -d
