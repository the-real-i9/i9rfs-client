package rfsSession

import (
	"context"
	"fmt"
	"i9rfs/client/appTypes"
	"i9rfs/client/helpers"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func authChallenge(authJwt string) error {
	if authJwt == "" {
		return fmt.Errorf("authentication required: please, login or signup")
	}

	connStream, err := helpers.WSConnect("ws://localhost:8000/api/app/session_user", authJwt)
	if err != nil {
		return fmt.Errorf("authorization: wsconn error: %s", err)
	}

	defer connStream.CloseNow()

	var recvData appTypes.WSResp
	// read response from connStream
	if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
		return fmt.Errorf("authorization: read error: %s", err)
	}

	if recvData.StatusCode != 200 {
		return fmt.Errorf("authentication error: %s", recvData.Error)
	}

	connStream.Close(websocket.StatusNormalClosure, "i am authorized")

	return nil
}
