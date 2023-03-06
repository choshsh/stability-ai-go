export APP_NAME=app
export BIN_DIR=bin

.PHONY: build
build: clean docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-tags lambda.norpc -ldflags '-s -w' \
		-o ${BIN_DIR}/${APP_NAME} main.go

.PHONY: serve
serve: docs
	cd cmd && go run main.go

.PHONY: deploy-dev
deploy-dev: build
	sls deploy function -f app -s dev

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: docs
docs:
	~/go/bin/swag init --parseDependency --parseInternal

