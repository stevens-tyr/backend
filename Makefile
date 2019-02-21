GO = go
GO_FILES = $(shell find . -maxdepth 1 -type f -name '*.go')
FMT = $(GO)imports -w 
LINT = $(GO)lint
TEST = $(GO) test
VET = $(GO) vet -composites=false
BUILD = $(GO) build
LIVE = gin run main.go

.PHONY: all test

build:
	$(BUILD) -o plague_doctor
live:
	env GIN_PORT=5000 BIN_APP_PORT=5555 $(LIVE) 
fmt:
	$(FMT) $(GO_FILES)
	$(FMT) $(shell find errors -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find middleware -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find models -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find models/cmsmodels -maxdepth 2 -type f -name '*.go')
	$(FMT) $(shell find models/usermodels -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find utils -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find api -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find api/cms -maxdepth 1 -type f -name '*.go')
	$(FMT) $(shell find api/auth -maxdepth 1 -type f -name '*.go')
lint:
	$(VET) $(GO_FILES)
	# $(VET) $(shell find models -maxdepth 1 -type f -name '*.go')
	# $(VET) $(shell find models/cmsmodels -maxdepth 2 -type f -name '*.go')
	# $(VET) $(shell find models/usermodels -maxdepth 1 -type f -name '*.go')
	# $(VET) $(shell find utils -maxdepth 1 -type f -name '*.go')
	$(VET) $(shell find api -maxdepth 1 -type f -name '*.go')
	$(VET) $(shell find api/cms -maxdepth 1 -type f -name '*.go')
	$(VET) $(shell find api/auth -maxdepth 1 -type f -name '*.go')
	$(LINT) $(GO_FILES)
	$(LINT) $(shell find models -maxdepth 1 -type f -name '*.go')
	$(LINT) $(shell find models/cmsmodels -maxdepth 2 -type f -name '*.go')
	$(LINT) $(shell find models/usermodels -maxdepth 1 -type f -name '*.go')
	$(LINT) $(shell find utils -maxdepth 1 -type f -name '*.go')
	$(LINT) $(shell find api -maxdepth 1 -type f -name '*.go')
	$(LINT) $(shell find api/cms -maxdepth 1 -type f -name '*.go')
	$(LINT) $(shell find api/auth -maxdepth 1 -type f -name '*.go')
test:
	$(TEST)
clean:
	rm -f plague_doctor
	rm -f log.json
all: fmt lint test build
