PLUGINS := $(wildcard plugins/*)
SUBDIRS = $(PLUGINS) mock otto

all: test build

test:
	go test ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
