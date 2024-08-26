package rfsSession

import "fmt"

func printWorkDir() {
	fmt.Printf("/home/%s%s\n", user.Username, workPath)
}
