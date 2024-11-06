LOCAL_BIN := $(CURDIR)/bin

GOENV := GOPRIVATE="github.com/ekubyshin"
PATH := $(PATH):$(LOCAL_BIN)

GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG := v1.61.0
$(GOLANGCI_BIN):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) $(GOLANGCI_TAG)
	chmod +x $(GOLANGCI_BIN)

.PHONY: lint
lint: $(GOLANGCI_BIN)
	$(GOENV) $(GOLANGCI_BIN) run --fix ./...

.PHONY: test
test:
	$(GOENV) go test ./... -count=1

.PHONY: bench
bench:
	$(GOENV) go test -run none -bench=. -benchtime=1s -cpu 8 -benchmem