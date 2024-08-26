package rfsSession

import (
	"context"
	"fmt"
	"i9pkgs/i9types"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func bashCommand(command string, cmdArgs []string, serverWorkPath string, connStream *websocket.Conn) {
	sendData := map[string]any{
		"workPath": serverWorkPath,
		"command":  command,
		"cmdArgs":  cmdArgs,
	}

	if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: write error: %s", command, w_err))
		return
	}

	var recvData i9types.WSResp

	if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: read error: %s", command, r_err))
		return
	}

	if recvData.StatusCode != 200 {
		fmt.Println(strings.Trim(recvData.Error, " \n"))
		return
	}

	fmt.Println(strings.Trim(recvData.Body.(string), " \n"))
}
