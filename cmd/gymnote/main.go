package main

import (
	"log"

	"gymnote/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
