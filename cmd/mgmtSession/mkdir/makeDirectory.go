package mkdir

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

	if cmdArgsLen := len(cmdArgs); cmdArgsLen != 1 {
		return fmt.Errorf("%d arguments provided, 1 required", cmdArgsLen)
	}

	if strings.ContainsAny(cmdArgs[0], "<>:\\|?*") {
		return fmt.Errorf("argument contains invalid characters: anyOf(<>:\\|?*\")")
	}

	if arg := cmdArgs[0]; strings.HasPrefix(arg, "\"") && !strings.HasSuffix(arg, "\"") {
		return fmt.Errorf("missing ending quote for argument")
	}

	return nil
}

func Run(command string, cmdArgs []string, workPath string, connStream *websocket.Conn) {
	ctx := context.Background()

	err := validateCmdArgs(cmdArgs)
	if err != nil {
		log.Println(fmt.Errorf("mkdir: %s", err))
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
		fmt.Printf("mkdir: %s\n", recvData.Error)
		return
	}
}
