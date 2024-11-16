CMD	= otto

all: 
	go build -v -o $(CMD)

pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o $(CMD)

clean:
	rm -f $(CMD)
	go clean
