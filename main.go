package main

import (
	"log"

	svc "github.com/chmouel/github-app-manifest-svc/pkg"
)

func main() {
	err := svc.Server()
	if err != nil {
		log.Fatal(err)
	}
}
