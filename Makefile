NAME=sms-query
VERSION=1.2.7

default: build-cli build-server

clean:
	@git clean -X -f

mod-download:
	@go mod download

fmt:
	@go fmt ./pkg/... ./cmd/$(NAME)

bindata:
	@go get -u github.com/jteeuwen/go-bindata/...
	@~/go/bin/go-bindata -o ./cmd/$(NAME)/bindata.go config/...

build-server: bindata
	@CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME) ./cmd/$(NAME)/bindata.go ./cmd/$(NAME)/server.go

build-server-rpi: bindata
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME) ./cmd/$(NAME)/bindata.go ./cmd/$(NAME)/server.go

build-cli: bindata
	@CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME)-cli ./cmd/$(NAME)/bindata.go ./cmd/$(NAME)/cli.go

build-cli-rpi: bindata
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X main.Version=$(VERSION)" -o $(NAME)-cli ./cmd/$(NAME)/bindata.go ./cmd/$(NAME)/cli.go

test:
	@CGO_ENABLED=0 go test -v -cover ./pkg/...

run-server: build-server
	@./$(NAME) -logtostderr -bind=127.0.0.1:8080

run-cli: build-cli
	@./$(NAME)-cli -logtostderr -from="$(from)" -query="$(query)" -send=$(send)

deploy: clean mod-download test build-server-rpi
	@./deploy.sh
