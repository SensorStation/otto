SUBDIRS := examples
PIENV	= env GOOS=linux GOARCH=arm GOARM=7

all: test $(SUBDIRS)

test:
	go test -cover ./...

test-v:
	go test -cover -v ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
