CMD	= otto

all: 
	go build -v -o $(CMD)

pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o $(CMD)

zero:
	env GOOS=linux GOARCH=arm GOARM=6 go build -o $(CMD)

clean:
	rm -f $(CMD) *.log
	go clean
