package authServices

import (
	"context"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/appTypes"
	"i9rfs/client/helpers"

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

func RegisterUser(connStream *websocket.Conn, signupSessionJwt2 string) error {
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
			"step":         "three",
			"sessionToken": signupSessionJwt2,
			"data": map[string]any{
				"username": username,
				"password": password,
			},
		}

		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return fmt.Errorf("signup: registerUser: write error: %s", err)
		}

		var recvData appTypes.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return fmt.Errorf("signup: registerUser: read error: %s", err)
		}

		// if app_err, continue for loop and ask for the code again, else break
		if recvData.StatusCode != 200 {
			fmt.Println(recvData.Error)
			continue
		}

		// Received data body
		var rcvdb struct {
			Msg     string
			User    map[string]any
			AuthJwt string
		}

		helpers.ParseTo(recvData.Body, &rcvdb)

		// store user data and auth_jwt
		appGlobals.AppDataStore.SetItem("user", rcvdb.User)
		appGlobals.AppDataStore.SetItem("auth_jwt", rcvdb.AuthJwt)
		appGlobals.AppDataStore.SetItem("i9rfs_work_path", "/")
		appGlobals.AppDataStore.Save()

		fmt.Println(rcvdb.Msg)

		return nil
	}
}
