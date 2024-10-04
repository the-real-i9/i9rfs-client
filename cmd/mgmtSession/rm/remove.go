package rm

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"
	"log"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func validateCmdArgs(cmdArgs []string) error {
	if !strings.HasPrefix(cmdArgs[0], "-") {
		if cmdArgsLen := len(cmdArgs); cmdArgsLen != 1 {
			return fmt.Errorf("%d arguments provided, 1 required", cmdArgsLen)
		}
	} else {
		if flg := cmdArgs[0]; flg != "-r" {
			return fmt.Errorf("invalid flag -- '%s'", flg)
		}

		if flgArgsLen := len(cmdArgs[1:]); flgArgsLen != 1 {
			return fmt.Errorf("%d arguments provided for flag '-r', 1 required", flgArgsLen)
		}
	}

	return nil
}

func Run(command string, cmdArgs []string, workPath string, connStream *websocket.Conn) {
	ctx := context.Background()

	err := validateCmdArgs(cmdArgs)
	if err != nil {
		log.Println(fmt.Errorf("rm: %s", err))
		return
	}

	sendData := map[string]any{
		"workPath": workPath,
		"command":  command,
		"cmdArgs":  cmdArgs,
	}

	if w_err := wsjson.Write(ctx, connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfsSession: rm: write error: %s", w_err))
		return
	}

	var recvData appTypes.WSResp

	if r_err := wsjson.Read(ctx, connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfsSession: rm: read error: %s", r_err))
		return
	}

	if recvData.StatusCode != 200 {
		fmt.Printf("mkdir: %s\n", recvData.Error)
		return
	}
}
