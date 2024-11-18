package fileUpload

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"
	"log"
	"os"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func filepathToBinary(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func Run(command string, cmdArgs []string, workPath string, connStream *websocket.Conn) {
	if cmdArgsLen := len(cmdArgs); cmdArgsLen != 2 {
		fmt.Printf("error: upload: %d arguments provided, 2 required\n", cmdArgsLen)
		return
	}

	fileLocation := cmdArgs[0]
	filename := cmdArgs[1]

	fileData, err := filepathToBinary(fileLocation)
	if err != nil {
		fmt.Printf("error: upload: %s", err)
	}

	sendData := map[string]any{
		"workPath": workPath,
		"command":  command,
		"cmdArgs":  []string{string(fileData), filename},
	}

	if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: write error: %s", command, w_err))
		return
	}

	var recvData appTypes.WSResp

	if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfsSession: %s: read error: %s", command, r_err))
		return
	}

	if recvData.StatusCode != 200 {
		fmt.Printf("%s: %s\n", command, recvData.Error)
		return
	}

	fmt.Println(recvData.Body)
}
