package otto

import (
	"log"
	"os"
)

func init() {
	f, err := os.OpenFile("otto.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	l = log.New(f, "otto: ", log.Ldate|log.Ltime|log.Lshortfile)
}
