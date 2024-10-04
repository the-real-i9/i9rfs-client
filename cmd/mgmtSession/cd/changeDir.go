package cd

import (
	"context"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/appTypes"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func Run(command string, cmdArgs []string, workPath string, connStream *websocket.Conn) {
	ctx := context.Background()

	if cmdArgsLen := len(cmdArgs); cmdArgsLen != 1 {
		log.Println(fmt.Errorf("error: cd: %d arguments provided, 1 required", cmdArgsLen))
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
