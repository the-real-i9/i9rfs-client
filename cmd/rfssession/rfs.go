package rfssession

import (
	"bufio"
	"fmt"
	"i9pkgs/i9helpers"
	"i9pkgs/i9services"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"nhooyr.io/websocket"
)

var workPath = ""
var user struct {
	Username string
}

func Launch() {
	if err := iAmAuthorized(); err != nil {
		fmt.Println(err)
		fmt.Println("Please, Login or create an account.")
		return
	}

	i9services.LocalStorage.GetItem("i9rfs_work_path", &workPath)
	i9services.LocalStorage.GetItem("user", &user)

	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/app/rfs", "")
	if err != nil {
		log.Println(fmt.Errorf("rfssession: Launch: connection error: %s", err))
		return
	}

	defer connStream.CloseNow()

	userAcc := color.New(color.Bold, color.FgGreen).Sprintf("i9rfs@%s", user.Username)

fsin:
	for {
		wpth := color.New(color.Bold, color.FgBlue).Sprintf("~%s", workPath)

		fmt.Printf("%s:%s$ ", userAcc, wpth)

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
		case "ls", "dir", "mv", "cp", "mkdir", "rmdir", "rm",
			"man", "gzip", "gunzip", "tar", "cat", "clear":
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
