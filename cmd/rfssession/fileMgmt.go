package rfssession

import (
	"context"
	"fmt"
	"i9pkgs/i9types"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func fileMgmtCommand(command string, cmdArgs []string, serverWorkPath string, connStream *websocket.Conn) {
	sendData := map[string]any{
		"workPath": serverWorkPath,
		"command":  command,
		"cmdArgs":  cmdArgs,
	}

	if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfssession: %s: write error: %s", command, w_err))
		return
	}

	var recvData i9types.WSResp

	if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfssession: %s: read error: %s", command, r_err))
		return
	}

	if recvData.Status == "f" {
		fmt.Println(strings.TrimRight(recvData.Error, " \n"))
		return
	}

	fmt.Println(strings.TrimRight(recvData.Body.(string), " \n"))
}
