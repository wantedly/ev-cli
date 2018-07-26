NAME     := ev
VERSION  := 1.0.1
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

OS := darwin linux

DIST_DIRS := find * -type d -exec

.PHONY: build
build: vendor/*
	go build $(LDFLAGS) -o $(NAME) *.go

.PHONY: install
install:
	@go install $(LDFLAGS)

.PHONY: cross-build
cross-build:
	@for os in $(OS); do \
		GOOS=$$os GOARCH=amd64 go build $(LDFLAGS) -o dist/$$os-amd64/$(NAME); \
	done

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-v$(VERSION)-{}.tar.gz {} \; && \
	cd ..

.PHONY: publish
publish:
	mkdir -p dist/publish
	@for os in $(OS); do \
		cp ./dist/$(NAME)-v$(VERSION)-$$os-amd64.tar.gz ./dist/publish/; \
	done

.PHONY: release
release:
	git tag v$(VERSION)
	git push origin v$(VERSION)

vendor/*: Gopkg.lock
	@$(MAKE) dep
	@dep ensure -v

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif
	@:
