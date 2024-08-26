package cmdAuthLogin

import (
	"context"
	"fmt"
	"i9rfs/client/appGlobals"
	"i9rfs/client/appTypes"
	"i9rfs/client/cmd/rfsSession"
	"i9rfs/client/helpers"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func Execute() {
	connStream, err := helpers.WSConnect("ws://localhost:8000/api/auth/login", "")
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

		var recvData appTypes.WSResp
		// read response from connStream
		if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
			log.Printf("signup: registerUser: read error: %s\n", err)
			return
		}

		if recvData.StatusCode != 200 {
			fmt.Println(recvData.Error)
			continue
		}

		var rcvdb struct {
			Msg      string
			User     map[string]any
			Auth_jwt string
		}

		helpers.ParseTo(recvData.Body, &rcvdb)

		// store user data and auth_jwt
		appGlobals.AppDataStore.SetItem("user", rcvdb.User)
		appGlobals.AppDataStore.SetItem("auth_jwt", rcvdb.Auth_jwt)
		appGlobals.AppDataStore.Save()

		fmt.Println(rcvdb.Msg)

		break
	}

	connStream.Close(websocket.StatusNormalClosure, "Login success!")

	rfsSession.Launch()
}
