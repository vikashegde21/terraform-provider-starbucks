default: install

build:
	go build -o terraform-provider-starbucks

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/starbucks/starbucks/0.1.0/darwin_arm64
	mv terraform-provider-starbucks ~/.terraform.d/plugins/registry.terraform.io/starbucks/starbucks/0.1.0/darwin_arm64/

test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

clean:
	rm -f terraform-provider-starbucks
	rm -rf .terraform/
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate*

.PHONY: build install test fmt lint clean
