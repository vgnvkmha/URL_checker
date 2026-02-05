package main

import (
	"URL_checker/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
