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

	arg := cmdArgs[0]

	if strings.ContainsAny(arg, "<>:\\|?*") {
		return fmt.Errorf("argument contains invalid characters: anyOf(<>:\\|?*)")
	}

	if strings.HasPrefix(arg, "/") || strings.HasSuffix(arg, "/") {
		return fmt.Errorf("argument can neither begin nor end with /")
	}

	for _, seg := range strings.Split(arg, "/") {
		if seg == "." {
			return fmt.Errorf("a directory name can't be '.'")
		}
		if seg == ".." {
			return fmt.Errorf("a directory name can't be '..'")
		}
		if strings.HasPrefix(seg, "\"") && !strings.HasSuffix(seg, "\"") {
			return fmt.Errorf("dir %s has no ending quote", seg)
		}
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
		log.Println(fmt.Errorf("rfsSession: mkdir: write error: %s", w_err))
		return
	}

	var recvData appTypes.WSResp

	if r_err := wsjson.Read(ctx, connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfsSession: mkdir: read error: %s", r_err))
		return
	}

	if recvData.StatusCode != 200 {
		fmt.Printf("mkdir: %s\n", recvData.Error)
		return
	}
}
