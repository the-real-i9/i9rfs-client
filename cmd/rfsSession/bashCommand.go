package rfsSession

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func bashCommand(command string, cmdArgs []string, serverWorkPath string, connStream *websocket.Conn) {
	ctx := context.Background()

	sendData := map[string]any{
		"workPath": serverWorkPath,
		"command":  command,
		"cmdArgs":  cmdArgs,
	}

	if w_err := wsjson.Write(ctx, connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: write error: %s", command, w_err))
		return
	}

	var recvData appTypes.WSResp

	if r_err := wsjson.Read(ctx, connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: read error: %s", command, r_err))
		return
	}

	if recvData.StatusCode != 200 {
		fmt.Println(strings.Trim(recvData.Error, " \n"))
		return
	}

	fmt.Println(strings.Trim(recvData.Body.(string), " \n"))
}
