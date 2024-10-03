package cmdAuthLogin

import (
	"fmt"
	"i9rfs/client/cmd/rfsSession"
	"i9rfs/client/helpers"
	"i9rfs/client/services/authServices"
	"log"

	"nhooyr.io/websocket"
)

func Execute() {
	connStream, err := helpers.WSConnect("ws://localhost:8000/api/auth/login", "")
	if err != nil {
		log.Printf("login: wsconn error: %s\n", err)
		return
	}

	defer connStream.CloseNow()

	err = authServices.Login(connStream)
	if err != nil {
		fmt.Println(err)
		return
	}

	connStream.Close(websocket.StatusNormalClosure, "Login success!")

	rfsSession.Launch()
}
