package rfssession

import (
	"context"
	"fmt"
	"i9pkgs/i9services"
	"i9pkgs/i9types"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func resolveToDestPath(workPath, destination string) (string, error) {
	if strings.HasPrefix(destination, "/") {
		return "", fmt.Errorf("invalid destination %s. Did you mean to prefix with ./ ?", destination)
	}
	dirs := strings.Split(destination, "/")

	destinationPath := workPath

	for _, dir := range dirs {
		if dir == "." {
			continue
		}

		if dir == ".." {
			if destinationPath == "" {
				// the user has specified an invalid directory,
				// one that possibly tries to go out of their user account directory
				return "", fmt.Errorf("no such file or directory")
			}

			// strip the last dir

			destinationPath = destinationPath[0:strings.LastIndex(destinationPath, "/")]
		} else {
			// append the dir
			destinationPath += "/" + dir
		}
	}

	return destinationPath, nil
}

func changeDirectory(cmdArgs []string, connStream *websocket.Conn) {
	if len(cmdArgs) > 1 {
		fmt.Println("error: cd: 2 arguments provided, 1 required")
		return
	}

	destinationPath, err := resolveToDestPath(workPath, cmdArgs[0])
	if err != nil {
		fmt.Printf("error: cd: %s\n", err)
		return
	}

	if destinationPath != "" {
		serverTestWorkPath := "/" + user.Username + destinationPath

		sendData := map[string]any{
			"workPath": serverTestWorkPath,
			"command":  "pex",
			"cmdArgs":  nil,
		}

		if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
			log.Println(fmt.Errorf("rfssession: command: cd: write error: %s", w_err))
			return
		}

		var recvData i9types.WSResp

		if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
			log.Println(fmt.Errorf("rfssession: command: cd: read error: %s", r_err))
			return
		}

		if recvData.Status == "f" {
			fmt.Printf("error: cd: %s\n", recvData.Error)
			return
		}

		if pathExists := recvData.Body.(bool); !pathExists {
			fmt.Println("error: cd: no such file or directory")
			return
		}
	}

	workPath = destinationPath

	i9services.LocalStorage.SetItem("i9rfs_work_path", workPath)

}
