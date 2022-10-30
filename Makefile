CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/alex1234562557/telegram-bot/cmd/bot

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test -count=1 ./...

prod:
	go run ${PACKAGE}

dev:
	go run ${PACKAGE} -develop

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/mocks/model/messages/messages_mocks.go
	${MOCKGEN} -source=internal/model/callbacks/incoming_clb.go -destination=internal/mocks/model/callbacks/callbacks_mocks.go
	${MOCKGEN} -source=internal/storage/storage.go -destination=internal/mocks/storage/storage_db_mocks.go
	${MOCKGEN} -source=internal/clients/rate/rateclient.go -destination=internal/mocks/clients/rate/rateclients_mocks.go
	${MOCKGEN} -source=internal/converter/converter.go -destination=internal/mocks/converter/converter_mocks.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	mkdir -p metrics/data
	sudo chmod -R 777 metrics/data
	sudo docker compose up

goose-status:
	goose -dir migrations  postgres "user=postgres password=pass dbname=postgres host=127.0.0.1 port=5432 sslmode=disable" status

goose-up:
	goose -dir migrations  postgres "user=postgres password=pass dbname=postgres host=127.0.0.1 port=5432 sslmode=disable" up

goose-reset:
	goose -dir migrations  postgres "user=postgres password=pass dbname=postgres host=127.0.0.1 port=5432 sslmode=disable" reset