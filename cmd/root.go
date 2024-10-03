package cmd

import (
	"fmt"
	"i9rfs/client/cmd/cmdAuth/cmdAuthLogin"
	"i9rfs/client/cmd/cmdAuth/cmdAuthLogout"
	"i9rfs/client/cmd/cmdAuth/cmdAuthSignup"
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
		// launch rfs session: this is the default action when no command is given
		rfsSession.Launch()

		return
	}

	// available commands include, authentication and help commands
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
