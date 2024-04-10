package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printHelp() {
	fmt.Println(
		`i9 Remote File System - CLI

Usage: i9rfs [command]

If command is empty, it launches the Remote File System with the existing user, if one exists.

The commands are:

   signup  - create an i9rfs account
   login   - login into i9rfs
   help    - print this usage guide
   logout  - logout of i9rfs
	 `)
}

func Execute() {
	args := os.Args[1:]

	if len(args) == 0 {
		// launch
		col := color.New(color.FgGreen, color.Italic, color.BlinkSlow)

		col.Println("Launching i9rfs...")

		return
	}

	switch args[0] {
	case "signup":
		// do signup
	case "login":
		// do login
	case "logout":
		// do logout
	case "help":
		// print help
		printHelp()
	default:
		printHelp()
	}
}
