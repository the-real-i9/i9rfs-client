package cmd

import (
	"fmt"
	"i9rfs/client/cmd/auth/cmdAuthLogin"
	"i9rfs/client/cmd/auth/cmdAuthLogout"
	"i9rfs/client/cmd/auth/cmdAuthSignup"
	"i9rfs/client/cmd/rfsSession"
	"os"
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
		rfsSession.Launch()

		return
	}

	switch args[0] {
	case "signup":
		cmdAuthSignup.Execute()
	case "login":
		cmdAuthLogin.Execute()
	case "logout":
		cmdAuthLogout.Execute()
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
