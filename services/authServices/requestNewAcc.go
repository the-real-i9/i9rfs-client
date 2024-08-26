package authServices

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"
	"i9rfs/client/helpers"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func RequestNewAccount(connStream *websocket.Conn) (signupSessionJwt string, newAccEmail string, rnaErr error) {
	for {
		// ask for email
		var email string

		fmt.Print("Email: ")
		fmt.Scanln(&email)

		sendData := map[string]any{"step": "one", "data": map[string]any{"email": email}}

		// send email data to WS server
		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return "", "", fmt.Errorf("signup: requestNewAccount: write error: %s", err)
		}

		var recvData appTypes.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return "", "", fmt.Errorf("signup: requestNewAccount: read error: %s", err)
		}

		var recvdb struct {
			SessionToken string
		}

		helpers.ParseTo(recvData.Body, &recvdb)

		// if app_err, continue for loop thereby asking for the email again, else break
		if recvData.StatusCode != 200 {
			fmt.Println(recvData.Error)
			continue
		}

		return "Bearer " + recvdb.SessionToken, email, nil
	}
}
