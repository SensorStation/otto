SUBDIRS := examples
PIENV	= env GOOS=linux GOARCH=arm GOARM=7

all: test $(SUBDIRS)

test:
	go test ./...

test-v:
	go test -v ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
