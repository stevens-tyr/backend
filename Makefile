GO = go
GO_FILES = $(wildcard *.go)
FMT = $(GO)imports -w 
LINT = $(GO)lint
TEST = $(GO) test
VET = $(GO) vet
BUILD = $(GO) build
LIVE = gin run main.go

build:
	$(BUILD) -o plague_doctor
live:
	export GIN_PORT=5000
	export BIN_APP_PORT=5001
	$(LIVE) 
fmt:
	$(FMT) $(GO_FILES)
lint:
	$(VET) $(GO_FILES)
	$(LINT) $(GO_FILES)
test:
	$(TEST)
clean:
	rm -f plague_doctor
	rm -f log.json
all: fmt lint test build
