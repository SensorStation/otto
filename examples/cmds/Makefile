all: 
	go build -v

pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build

clean:
	rm -f cmd *.log
	go clean
