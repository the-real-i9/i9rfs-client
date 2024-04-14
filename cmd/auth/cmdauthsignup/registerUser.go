package cmdauthsignup

import (
	"context"
	"fmt"
	"i9pkgs/i9helpers"
	"i9pkgs/i9services"
	"i9pkgs/i9types"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func validateUserInfo(username string, password string) error {
	if len(username) < 3 {
		return fmt.Errorf("username too short, must be minimum of 3 characters")
	}

	if len(password) < 8 {
		return fmt.Errorf("password too short, must be minimum of 8 characters")
	}

	return nil
}

func registerUser(connStream *websocket.Conn, signupSessionJwt string) error {
	for {
		// ask for username
		// ask for password
		var (
			username string
			password string
		)

		for {
			fmt.Print("Username: ")
			fmt.Scanln(&username)

			fmt.Print("Password: ")
			fmt.Scanln(&password)

			if err := validateUserInfo(username, password); err != nil {
				fmt.Println(err)
				continue
			}

			break
		}

		// write user data along with signup_session_jwt to WS server
		sendData := map[string]any{
			"step": "three",
			"data": map[string]any{
				"user_info": map[string]any{
					"username": username,
					"password": password,
				},
				"signup_session_jwt": signupSessionJwt,
			},
		}

		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return fmt.Errorf("signup: registerUser: write error: %s", err)
		}

		var recvData i9types.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return fmt.Errorf("signup: registerUser: read error: %s", err)
		}

		// if app_err, continue for loop and ask for the code again, else break
		if recvData.Status == "f" {
			fmt.Println(recvData.Error)
			continue
		}

		var rcvdb struct {
			Msg      string
			User     map[string]any
			Auth_jwt string
		}

		i9helpers.ParseTo(recvData.Body, &rcvdb)

		// store user data and auth_jwt
		i9services.LocalStorage.SetItem("user", rcvdb.User, "auth_jwt", rcvdb.Auth_jwt)

		fmt.Println(rcvdb.Msg)

		return nil
	}
}
