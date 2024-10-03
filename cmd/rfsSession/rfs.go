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

	if authJwt == "" {
		log.Println("authentication required: please, login or signup")
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

		switch command {
		case "cd":
			changeDirectory(command, cmdArgs, workPath, connStream)
		case "upload", "up":
			uploadFile(command, cmdArgs, workPath, connStream)
		case "download", "down":
			downloadFile(command, cmdArgs, workPath, connStream)
		case "ls", "dir", "mv", "cp", "mkdir", "rmdir", "rm", "clear":
			bashCommand(command, cmdArgs, workPath, connStream)
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
