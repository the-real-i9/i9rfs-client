package cmd

import (
	"fmt"
	"i9rfs/client/cmd/authSession/authCmdLogin"
	"i9rfs/client/cmd/authSession/authCmdLogout"
	"i9rfs/client/cmd/authSession/authCmdSignup"
	"i9rfs/client/cmd/mgmtSession"
	"os"
)

func printHelp() {
	fmt.Println(
		`i9 Remote File System - CLI

Usage: i9rfs [command]

If command is empty, it begines the Remote File System with the existing user, if one exists.

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
		// begin management session: this is the default action when no command is given
		mgmtSession.Begin()

		return
	}

	// available commands include, authentication and help commands
	switch args[0] {
	case "signup":
		authCmdSignup.Execute()
	case "login":
		authCmdLogin.Execute()
	case "logout":
		authCmdLogout.Execute()
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
