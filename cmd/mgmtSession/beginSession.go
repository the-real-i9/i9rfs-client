package mgmtSession

import (
	"bufio"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/cmd/mgmtSession/cd"
	"i9rfs/client/cmd/mgmtSession/mkdir"
	"i9rfs/client/cmd/mgmtSession/rm"
	"i9rfs/client/cmd/mgmtSession/rmdir"
	"i9rfs/client/helpers"
	"log"
	"os"
	"strings"
	"unicode"

	"nhooyr.io/websocket"
)

var workPath = ""
var user struct {
	Username string
}

func Begin() {
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
		log.Println(fmt.Errorf("mgmtSession: Begin: connection error: %s", err))
		return
	}

	defer connStream.CloseNow()

	comp := fmt.Sprintf("%s@i9frs", user.Username)

fsin:
	for {
		prompt := fmt.Sprintf("%s:%s$ ", comp, fmt.Sprintf("~%s", strings.TrimSuffix(workPath, "/")))

		fmt.Print(prompt)

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		// we want to consider quoted strings contanining whitespace as one
		inQuote := false
		cmdLine := strings.FieldsFunc(strings.TrimSpace(input.Text()), func(r rune) bool {
			if r == '"' {
				if !inQuote {
					inQuote = true
				} else {
					inQuote = false
				}
			}

			if (unicode.IsSpace(r) && inQuote) || !unicode.IsSpace(r) {
				return false
			}

			return true
		})

		command := cmdLine[0]
		cmdArgs := cmdLine[1:]

		switch command {
		case "cd":
			cd.Run(command, cmdArgs, &workPath, connStream)
		case "upload", "up":
			uploadFile(command, cmdArgs, workPath, connStream)
		case "download", "down":
			downloadFile(command, cmdArgs, workPath, connStream)
		case "mkdir":
			mkdir.Run(command, cmdArgs, workPath, connStream)
		case "rmdir":
			rmdir.Run(command, cmdArgs, workPath, connStream)
		case "rm":
			rm.Run(command, cmdArgs, workPath, connStream)
		case "ls", "dir", "mv", "cp", "clear":
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
