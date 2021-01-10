PKGNAME  = abi-wizard
DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

GOPROJROOT  = $(GOSRC)/$(PROJREPO)

GOLDFLAGS   = -ldflags "-s -w"
GOTAGS      = ""
GOCC        = go
GOFMT       = $(GOCC) fmt -x
GOGET       = $(GOCC) get $(GOLDFLAGS)
GOBUILD     = $(GOCC) build -v $(GOLDFLAGS) $(GOTAGS)
GOTEST      = $(GOCC) test
GOVET       = $(GOCC) vet
GOINSTALL   = $(GOCC) install $(GOLDFLAGS)
GOBUILDDEP  = $(GOINSTALL)
GOCLEANDEP  = $(GOCC) clean -cache -modcache

include Makefile.waterlog

GOLINT    = golint -set_exit_status

all: build

build:
	@$(call stage,BUILD)
	@$(GOBUILD)
	@$(call pass,BUILD)

test:
	@$(call stage,TEST)
	@$(GOTEST) ./...
	@$(call pass,TEST)

validate:
	@$(call stage,FORMAT)
	@$(GOFMT) ./...
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(call task,Running 'go vet'...)
	@$(GOVET) ./...
	@$(call pass,VET)
	@$(call stage,LINT)
	@$(call task,Running 'golint'...)
	@$(GOLINT) `go list ./... | grep -v vendor`
	@$(call pass,LINT)

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing executable...)
	@rm $(PKGNAME)
	@$(call pass,CLEAN)
