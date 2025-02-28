HOSTNAME=registry.terraform.io
NAMESPACE=amnic
NAME=cortex
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_arm64

default: install

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: release
release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test:
	go test ./... -timeout=30s -parallel=4 -count=1 -v

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -timeout 120m -v

.PHONY: clean
clean:
	rm -f examples/terraform.tfstate
	rm -f examples/terraform.tfstate.backup

dev.tfrc:
	echo 'provider_installation {' >> dev.tfrc
	echo '  dev_overrides {' >> dev.tfrc
	echo '    "form3tech-oss/cortex" = "$(CURDIR)"' >> dev.tfrc
	echo '  }' >> dev.tfrc
	echo '  direct {}' >> dev.tfrc
	echo '}' >> dev.tfrc

.PHONY: cortex-up
cortex-up:
	docker-compose up -d

.PHONY: cortex-down
cortex-down:
	docker-compose down
