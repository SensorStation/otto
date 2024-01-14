PLUGINS := $(wildcard plugins/*)
SUBDIRS = $(PLUGINS) otto

all: test $(SUBDIRS)

test:
	go test ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
