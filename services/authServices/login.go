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

func Login(connStream *websocket.Conn) error {
	var (
		emailOrUsername string
		password        string
	)

	for {
		fmt.Print("Email/Username: ")
		fmt.Scanln(&emailOrUsername)

		fmt.Print("Password: ")
		fmt.Scanln(&password)

		sendData := map[string]any{
			"emailOrUsername": emailOrUsername,
			"password":        password,
		}

		if err := wsjson.Write(context.Background(), connStream, sendData); err != nil {
			return fmt.Errorf("login: write error: %s\n", err)
		}

		var recvData appTypes.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			return fmt.Errorf("signup: registerUser: read error: %s\n", err)
		}

		if recvData.StatusCode != 200 {
			fmt.Println("\n", recvData.Error)
			continue
		}

		var rcvdb struct {
			Msg     string
			User    map[string]any
			AuthJwt string
		}

		helpers.ParseTo(recvData.Body, &rcvdb)

		// store user data and auth_jwt
		appGlobals.AppDataStore.SetItem("user", rcvdb.User)
		appGlobals.AppDataStore.SetItem("auth_jwt", rcvdb.AuthJwt)
		appGlobals.AppDataStore.Save()

		fmt.Printf("%s\n\n", rcvdb.Msg)

		return nil
	}
}
