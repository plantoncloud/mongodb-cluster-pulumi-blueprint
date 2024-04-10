.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf ${build_dir}

.PHONY: build
build: clean deps vet fmt test

.PHONY: test
test:
	go test -v ./...

.PHONY: update-deps
update-deps:
	go get github.com/plantoncloud/planton-cloud-apis@latest
	go get github.com/plantoncloud/pulumi-stack-runner-go-sdk
