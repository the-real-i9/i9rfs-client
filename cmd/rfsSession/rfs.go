package rfsSession

import (
	"bufio"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/helpers"
	"log"
	"os"
	"strings"

	"nhooyr.io/websocket"
)

var workPath = ""
var user struct {
	Username string
}

func Launch() {
	var authJwt string

	appGlobals.AppDataStore.GetItem("auth_jwt", &authJwt)

	if err := authChallenge(authJwt); err != nil {
		fmt.Println(err)
		return
	}

	appGlobals.AppDataStore.GetItem("i9rfs_work_path", &workPath)
	appGlobals.AppDataStore.GetItem("user", &user)

	connStream, err := helpers.WSConnect("ws://localhost:8000/api/app/rfs", authJwt)
	if err != nil {
		log.Println(fmt.Errorf("rfsSession: Launch: connection error: %s", err))
		return
	}

	defer connStream.CloseNow()

	comp := fmt.Sprintf("%s@i9frs", user.Username)

fsin:
	for {
		prompt := fmt.Sprintf("%s:%s$ ", comp, fmt.Sprintf("~%s", workPath))

		fmt.Print(prompt)

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		cmdLine := strings.Split(input.Text(), " ")
		command := cmdLine[0]
		cmdArgs := cmdLine[1:]

		serverWorkPath := "/" + user.Username + workPath

		switch command {
		case "cd":
			changeDirectory(cmdArgs, connStream)
		case "pwd":
			printWorkDir()
		case "upload", "up":
			uploadFile(command, cmdArgs, serverWorkPath, connStream)
		case "download", "down":
			downloadFile(command, cmdArgs, serverWorkPath, connStream)
		case "ls", "dir", "mv", "cp", "mkdir", "rmdir", "rm", "clear":
			bashCommand(command, cmdArgs, serverWorkPath, connStream)
		case "exit":
			fmt.Println("exiting...")
			break fsin
		case "":
			continue fsin
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}

	connStream.Close(websocket.StatusNormalClosure, "exiting...")
}
