pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -v . 

test:
	go test ./...

test-v:
	go test -v ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all test build $(SUBDIRS)
