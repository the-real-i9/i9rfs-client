package cmdauthsignup

import (
	"context"
	"fmt"
	"i9pkgs/i9types"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func verifyEmail(connStream *websocket.Conn, signupSessionJwt string, newAccEmail string) error {
	for {
		// ask for verf code
		var code int

		fmt.Println("Note: If you have issues receiving code, press `Crtl + C` to terminate the session and start signup over.")
		for {
			fmt.Printf("Enter the 6-digit code sent to '%s':\n", newAccEmail)

			if _, err := fmt.Scanf("%d", &code); err != nil {
				fmt.Println(err)
				continue
			}

			break
		}

		// send code along with signup_session_jwt to WS server
		sendData := map[string]any{"step": "two", "data": map[string]any{"code": code, "signup_session_jwt": signupSessionJwt}}

		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return fmt.Errorf("signup: verifyEmail: write error: %s", err)
		}

		var recvData i9types.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return fmt.Errorf("signup: verifyEmail: read error: %s", err)
		}

		// if app_err, continue for loop and ask for the code again, else break
		if recvData.Status == "f" {
			fmt.Println(recvData.Error)
			continue
		}

		fmt.Println(recvData.Body)

		return nil
	}
}
