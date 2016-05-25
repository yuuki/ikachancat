BIN = droot

all: clean build

build: deps gen
	go build -o $(BIN) ./cmd

fmt: deps
	gofmt -s -w .

validate: lint
	go vet ./...
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

lint:
	out="$$(golint ./...)"; \
	if [ -n "$$(golint ./...)" ]; then \
		echo "$$out"; \
		exit 1; \
	fi

patch: gobump
	./script/release.sh patch

minor: gobump
	./script/release.sh minor

gobump:
	go get github.com/motemen/gobump/cmd/gobump

deps:
	go get -d -v ./...

clean:
	go clean

.PHONY: test build lint patch minor gobump deps clean
