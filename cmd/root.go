package cmd

import (
	"fmt"
	"i9rfs/client/cmd/auth/cmdauthlogin"
	"i9rfs/client/cmd/auth/cmdauthlogout"
	"i9rfs/client/cmd/auth/cmdauthsignup"
	"i9rfs/client/cmd/rfssession"
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
		rfssession.Launch()

		return
	}

	switch args[0] {
	case "signup":
		cmdauthsignup.Execute()
	case "login":
		cmdauthlogin.Execute()
	case "logout":
		cmdauthlogout.Execute()
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
