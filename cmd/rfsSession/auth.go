package rfsSession

import (
	"context"
	"fmt"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func authChallenge() error {
	var authJwt string
	i9services.LocalStorage.GetItem("auth_jwt", &authJwt)

	if authJwt == "" {
		return fmt.Errorf("authentication required: please, login or create an account")
	}

	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/app/get_session_user", authJwt)
	if err != nil {
		return fmt.Errorf("authorization: wsconn error: %s", err)
	}

	defer connStream.CloseNow()

	var recvData i9types.WSResp
	// read response from connStream
	if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
		return fmt.Errorf("authorization: read error: %s", err)
	}

	if recvData.StatusCode != 200 {
		return fmt.Errorf("authentication required: please, login or create an account")
	}

	connStream.Close(websocket.StatusNormalClosure, "i am authorized")

	return nil
}
