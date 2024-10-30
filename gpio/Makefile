test:
	go test ./...

test-v:
	go test -v ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
