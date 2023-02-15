.PHONY: all build run go_tool clean help

BUILD_DIR=build
BUILD_ARGS=-ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
VERSION=$(shell git describe --tags)
BINARY_NAME=biliDownload
BINARY_LINUX_AMD64=${BINARY_NAME}-linux-amd64-${VERSION}
BINARY_LINUX_ARM64=${BINARY_NAME}-linux-arm64-${VERSION}
BINARY_MACOS_AMD64=${BINARY_NAME}-mac-amd64-${VERSION}
BINARY_MACOS_ARM64=${BINARY_NAME}-mac-arm64-${VERSION}
BINARY_WINDOWS_AMD64=${BINARY_NAME}-win-amd64-${VERSION}.exe
BINARY_WINDOWS_ARM64=${BINARY_NAME}-win-arm64-${VERSION}.exe
BINARY_PACK=${BINARY_NAME}-${VERSION}.zip

all: help

build: go_tool clean
	mkdir ${BUILD_DIR}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_LINUX_AMD64}
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_LINUX_ARM64}
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_MACOS_AMD64}
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_MACOS_ARM64}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_WINDOWS_AMD64}
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build ${BUILD_ARGS} -o ${BUILD_DIR}/${BINARY_WINDOWS_ARM64}
	cd ${BUILD_DIR} && 7zz a ${BINARY_PACK} *amd64* *arm64*

run:
	@go run ./

go_tool:
	go fmt ./
	go vet ./

clean:
	@if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR} ; fi

help:
	@echo "make - 查看帮助信息"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make go_tool - 运行 Go 工具 'fmt' and 'vet'"