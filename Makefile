.PHONY: clean deps format check test build build-linux push push-remote setup container

all: clean deps format check build

jenkins: clean deps format build-linux

build: test
	@echo "==> Building..."
	CGO_ENABLED=0 ${GO} build -o ${APP_NAME}

build-linux: test
	@echo "==> Building..."
	@echo "Local build" > .BUILD.txt
	@echo `hostname` >> .BUILD.txt
	@echo `date` >> .BUILD.txt
	cat .BUILD.txt
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -ldflags="-w -s" -o ${APP_NAME}-linux

container: build-linux
	docker build -t ${REMOTE_IMAGE_REPO}/${REPO_ORG}/kaas-entitlement-controller:${VERSION} .

push-remote:
	docker push ${REMOTE_IMAGE_REPO}/${REPO_ORG}/kaas-entitlement-controller:${VERSION}

push: container push-remote
	@echo "Done"

deploy:
	oc apply -f ./rbac
	oc apply -f ./deployment.yaml -n ${NAMESPACE}

clean:
	@echo "==> Cleaning..."
	rm -f coverage.out report.json
	rm -f ${APP_NAME}
	rm -f ${APP_NAME}-linux

deps: 
	@echo "==> Getting Dependencies..."
	${GO} mod tidy
	${GO} mod download

test: 
	@echo "==> Testing..."
	CGO_ENABLED=0 ${GO} test -tags test -v -covermode=atomic -count=1 ./... -coverprofile coverage.out
	${GO} test -race -tags test -covermode=atomic -count=1 ./... -json > report.json
	${GO} tool cover -func=coverage.out

format:
	@echo "==> Code Formatting..."
	${GO} fmt  ./pkg/...

check: format
	@echo "==> Code Check..."
	golangci-lint run -c .golangci.yml

setup:
	@echo "==> Setup..."
	${GO} get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.14.0
