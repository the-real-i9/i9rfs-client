package main

import (
	"i9pkgs/i9services"
	"i9rfs/client/cmd"
	"log"
)

func main() {
	if err := i9services.InitLocalStorage("i9rfs", "localStorage.json"); err != nil {
		log.Fatal(err)
	}

	i9services.LocalStorage.SetItem("bar", "foo")

	cmd.Execute()
}
