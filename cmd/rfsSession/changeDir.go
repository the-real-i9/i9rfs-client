package rfsSession

import (
	"context"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/appTypes"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func changeToTargetPath(workPath, targetPath string) (string, error) {
	if strings.HasPrefix(targetPath, "/") {
		return "", fmt.Errorf("invalid target path %s. Did you mean to prefix with ./ instead?", targetPath)
	}
	dirs := strings.Split(targetPath, "/")

	newWorkPath := workPath

	for _, dir := range dirs {
		if dir == "." {
			continue
		}

		if dir == ".." {
			if newWorkPath == "" {
				// the user has specified an invalid directory,
				// one that possibly tries to go out of their user account directory
				return "", fmt.Errorf("no such file or directory")
			}

			// strip the last dir
			newWorkPath = newWorkPath[0:strings.LastIndex(newWorkPath, "/")]
		} else {
			// append the dir
			newWorkPath += "/" + dir
		}
	}

	return newWorkPath, nil
}

func changeDirectory(cmdArgs []string, connStream *websocket.Conn) {
	ctx := context.Background()

	if cmdArgsLen := len(cmdArgs); cmdArgsLen > 1 {
		fmt.Printf("error: cd: %d arguments provided, 1 required\n", cmdArgsLen)
		return
	}

	newWorkPath, err := changeToTargetPath(workPath, cmdArgs[0])
	if err != nil {
		fmt.Printf("cd: %s\n", err)
		return
	}

	if newWorkPath != "" {
		serverTestWorkPath := "/" + user.Username + newWorkPath

		sendData := map[string]any{
			"workPath": serverTestWorkPath,
			"command":  "pex", // "path exist"ence challenge
			"cmdArgs":  nil,
		}

		if w_err := wsjson.Write(ctx, connStream, sendData); w_err != nil {
			log.Println(fmt.Errorf("rfsSession: cd: write error: %s", w_err))
			return
		}

		var recvData appTypes.WSResp

		if r_err := wsjson.Read(ctx, connStream, &recvData); r_err != nil {
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

	workPath = newWorkPath

	appGlobals.AppDataStore.SetItem("i9rfs_work_path", workPath)
	appGlobals.AppDataStore.Save()

}
