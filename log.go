package otto

import (
	"log"
	"os"
)

func init() {
	f, err := os.OpenFile("otto.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("This is a test log entry")
}
