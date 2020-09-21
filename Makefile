run:
	docker-compose -f docker-compose.yml up --build -d

stop:
	docker-compose -f docker-compose.yml down

lint:
	GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	cd jobber/ ; \
	golangci-lint run

unit_tests:
	cd jobber/ ; \
	go test ./...

integration_tests:
	docker-compose -f docker-compose-test.yml up -V --abort-on-container-exit --exit-code-from testserver
	docker-compose -f docker-compose-test.yml down

build:
	cd jobber/ ; \
	go build cmd/main.go
