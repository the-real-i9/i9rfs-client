package authServices

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func VerifyEmail(connStream *websocket.Conn, signupSessionJwt string, newAccEmail string) (signupSessionJwt2 string, veError error) {
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
		sendData := map[string]any{"step": "two", "sessionToken": signupSessionJwt, "data": map[string]any{"code": code}}

		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return "", fmt.Errorf("signup: verifyEmail: write error: %s", err)
		}

		var recvData appTypes.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return "", fmt.Errorf("signup: verifyEmail: read error: %s", err)
		}

		// if app_err, continue for loop and ask for the code again, else break
		if recvData.StatusCode != 200 {
			fmt.Println(recvData.Error)
			continue
		}

		return recvData.Body.(string), nil
	}
}
