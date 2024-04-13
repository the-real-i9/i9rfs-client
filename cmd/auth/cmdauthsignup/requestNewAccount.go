package cmdauthsignup

import (
	"context"
	"fmt"
	"i9pkgs/i9services"
	"i9pkgs/i9types"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func requestNewAccount(connStream *websocket.Conn) error {
	for {
		// ask for email
		var email string

		fmt.Print("Email: ")
		fmt.Scanln(&email)

		sendData := map[string]any{"step": "one", "data": email}

		// send email data to WS server
		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return fmt.Errorf("signup: requestNewAccount: write error: %s", err)
		}

		var recvData i9types.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return fmt.Errorf("signup: requestNewAccount: read error: %s", err)
		}

		// if app_err, continue for loop thereby asking for the email again, else break
		if recvData.Status == "f" {
			fmt.Println(recvData.Error)
			continue
		}

		i9services.LocalStorage.SetItem("signup_session_jwt", recvData.Body)

		break

	}

	return nil
}
