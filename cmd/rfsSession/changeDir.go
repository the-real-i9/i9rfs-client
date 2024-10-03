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

func changeToTargetPath(currentWorkPath, targetPath string) (string, error) {
	if strings.HasPrefix(targetPath, "/") {
		return "", fmt.Errorf("invalid target path %s. Did you mean to prefix with ./ instead?", targetPath)
	}
	dirs := strings.Split(targetPath, "/")

	newWorkPath := currentWorkPath

	for _, dir := range dirs {
		if dir == "." {
			continue
		} else if dir == ".." {
			if newWorkPath == "/" {
				// the user has specified an invalid directory,
				// one that possibly tries to go out of their user account directory
				continue
			}

			// strip the last dir
			// if newWorkPath was the last directory in the root
			// the code line below will make it an empty string
			newWorkPath = newWorkPath[0:strings.LastIndex(newWorkPath, "/")]
			// hence, we check and restore to root
			if newWorkPath == "" {
				newWorkPath = "/"
			}
		} else {
			// append the dir
			if newWorkPath == "/" { // root
				newWorkPath += dir
			} else { // non-root
				newWorkPath += "/" + dir
			}
		}
	}

	return newWorkPath, nil
}

func changeDirectory(command string, cmdArgs []string, workPath string, connStream *websocket.Conn) {
	ctx := context.Background()

	if cmdArgsLen := len(cmdArgs); cmdArgsLen > 1 {
		fmt.Printf("error: cd: %d arguments provided, 1 required\n", cmdArgsLen)
		return
	}

	sendData := map[string]any{
		"workPath": workPath,
		"command":  command,
		"cmdArgs":  cmdArgs,
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

	newWorkPath := recvData.Body.(string)

	workPath = newWorkPath

	appGlobals.AppDataStore.SetItem("i9rfs_work_path", workPath)
	appGlobals.AppDataStore.Save()

}
