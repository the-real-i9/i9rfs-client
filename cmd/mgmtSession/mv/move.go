package mv

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
	if cmdArgsLen := len(cmdArgs); cmdArgsLen != 2 {
		return fmt.Errorf("%d arguments provided, 2 required", cmdArgsLen)
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
		sourceSegments := strings.Split(cmdArgs[0], "/")

		errMsgRepl := strings.NewReplacer("$dest/$source_last_seg", fmt.Sprintf("'%s/%s'", cmdArgs[1], sourceSegments[len(sourceSegments)-1]), "$source", fmt.Sprintf("'%s'", cmdArgs[0]), "$dest", fmt.Sprintf("'%s'", cmdArgs[1]))

		fmt.Printf("mv: %s\n", errMsgRepl.Replace(recvData.Error))
		return
	}
}
