package main

import (
	"i9rfs/client/cmd"
	"i9rfs/client/initializers"
	"log"
)

func main() {
	if err := initializers.InitApp(); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
