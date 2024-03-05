GO=go1.22.0
BUILD=build

generate:
	${GO} generate ./...

build: generate
	mkdir -p ${BUILD}
	${GO} build -o ${BUILD} ./...

run: build
	./${BUILD}/cmd
