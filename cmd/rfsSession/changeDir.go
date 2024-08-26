package rfsSession

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
	if cmdArgsLen := len(cmdArgs); cmdArgsLen > 1 {
		fmt.Printf("error: cd: %d arguments provided, 1 required\n", cmdArgsLen)
		return
	}

	destinationPath, err := resolveToDestPath(workPath, cmdArgs[0])
	if err != nil {
		fmt.Printf("cd: %s\n", err)
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
			log.Println(fmt.Errorf("rfsSession: cd: write error: %s", w_err))
			return
		}

		var recvData i9types.WSResp

		if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
			log.Println(fmt.Errorf("rfsSession: cd: read error: %s", r_err))
			return
		}

		if recvData.StatusCode != 200 {
			fmt.Printf("cd: %s\n", recvData.Error)
			return
		}

		if pathExists := recvData.Body.(bool); !pathExists {
			fmt.Println("cd: no such file or directory")
			return
		}
	}

	workPath = destinationPath

	i9services.LocalStorage.SetItem("i9rfs_work_path", workPath)

}
