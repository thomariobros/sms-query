NAME=sms-query
VERSION=1.2.7

default: build-cli build-server test

clean:
	@git clean -X -f

mod-download:
	@go mod download

fmt:
	@gofmt -s -l -w .

fmt-check:
	@./scripts/fmt-check.sh

bindata:
	@go get -u github.com/jteeuwen/go-bindata/...
	@~/go/bin/go-bindata -o ./pkg/run/bindata.go -pkg run config/i18n/...

build-server: bindata
	@CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME) ./cmd/$(NAME)/server.go

build-cli: bindata
	@CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME)-cli ./cmd/$(NAME)/cli.go

test:
	@CGO_ENABLED=0 go test -v -cover ./pkg/...

run-server: build-server
	@./$(NAME) -logtostderr -config=config/config.yml -bind=127.0.0.1:8080

run-cli: build-cli
	@./$(NAME)-cli -logtostderr -config=config/config.yml -from="$(from)" -query="$(query)" -send=$(send)
