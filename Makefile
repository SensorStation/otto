SUBDIRS := examples
PIENV	= env GOOS=linux GOARCH=arm GOARM=7

all: test $(SUBDIRS)

init:
	git update --init 

test:
	go test -cover ./...

test-v:
	go test -cover -v ./...

html:
	go test -coverprofile=/home/rusty/cover.out ./...
	go tool cover -html=/home/rusty/cover.out

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
