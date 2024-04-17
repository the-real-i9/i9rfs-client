package cmdauthlogin

import (
	"context"
	"fmt"
	"i9pkgs/i9helpers"
	"i9pkgs/i9services"
	"i9pkgs/i9types"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func Execute() {
	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/auth/login", "")
	if err != nil {
		log.Printf("login: wsconn error: %s\n", err)
		return
	}

	defer connStream.CloseNow()

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
			log.Printf("login: write error: %s\n", err)
		}

		var recvData i9types.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			log.Printf("signup: registerUser: read error: %s\n", err)
			return
		}

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

		break
	}

	connStream.Close(websocket.StatusNormalClosure, "Login success!")
}
