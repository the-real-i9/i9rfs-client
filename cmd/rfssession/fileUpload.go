package rfssession

import (
	"context"
	"fmt"
	"i9pkgs/i9types"
	"log"
	"os"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func filepathToBinary(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func uploadFile(cmdArgs []string, connStream *websocket.Conn) {
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
		"command":  "upload",
		"cmdArgs":  []string{string(fileData), filename},
	}

	if w_err := wsjson.Write(context.Background(), connStream, sendData); w_err != nil {
		log.Println(fmt.Errorf("rfssession: upload: write error: %s", w_err))
		return
	}

	var recvData i9types.WSResp

	if r_err := wsjson.Read(context.Background(), connStream, &recvData); r_err != nil {
		log.Println(fmt.Errorf("rfssession: upload: read error: %s", r_err))
		return
	}

	if recvData.Status == "f" {
		fmt.Printf("error: upload: %s\n", recvData.Error)
		return
	}

	fmt.Print(recvData.Body)
}
