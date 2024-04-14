package rfsinter

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

func resolveToTestWorkPath(workPath, target string) (string, error) {
	if strings.HasPrefix(target, "/") {
		return "", fmt.Errorf("invalid target %s. Did you mean to prefix with ./ ?", target)
	}
	dirs := strings.Split(target, "/")

	testWorkPath := workPath

	for _, dir := range dirs {
		if dir == "." {
			continue
		}

		if dir == ".." {
			if testWorkPath == "" {
				// the user has specified an invalid directory,
				// one that possibly tries to go out of their user account directory
				return "", fmt.Errorf("no such file or directory")
			}

			// strip the last dir

			testWorkPath = testWorkPath[0:strings.LastIndex(testWorkPath, "/")]
		} else {
			// append the dir
			testWorkPath += "/" + dir
		}
	}

	return testWorkPath, nil
}

func changeDirectory(cmdArgs []string, connStream *websocket.Conn) {
	if len(cmdArgs) > 1 {
		fmt.Println("error: cd: 2 arguments provided, 1 required")
		return
	}

	testWorkPath, err := resolveToTestWorkPath(workPath, cmdArgs[0])
	if err != nil {
		fmt.Printf("error: cd: %s\n", err)
		return
	}

	if testWorkPath != "" {
		serverTestWorkPath := "/" + user.Username + testWorkPath

		sendData := map[string]any{
			"workPath": serverTestWorkPath,
			"command":  "cd",
			"cmdArgs":  nil,
		}

		if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
			log.Println(fmt.Errorf("rfsinter: Launch: command: cd: write error: %s", w_err))
			return
		}

		var recvData i9types.WSResp

		if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
			log.Println(fmt.Errorf("rfsinter: Launch: command: cd: read error: %s", r_err))
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

	workPath = testWorkPath

	i9services.LocalStorage.SetItem("i9rfs_work_path", workPath)

}
