package rfsinter

import (
	"context"
	"fmt"
	"i9pkgs/i9helpers"
	"i9pkgs/i9services"
	"i9pkgs/i9types"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func iAmAuthorized() error {
	var authJwt string
	i9services.LocalStorage.GetItem("auth_jwt", &authJwt)

	if authJwt == "" {
		return fmt.Errorf("error: user authentication required")
	}

	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/app/get_session_user", authJwt)
	if err != nil {
		return fmt.Errorf("authorization error: %s", err)
	}

	defer connStream.CloseNow()

	var recvData i9types.WSResp
	// read response from connStream
	if err := wsjson.Read(context.Background(), connStream, &recvData); err != nil {
		return fmt.Errorf("authorization: read error: %s", err)
	}

	if recvData.Status == "f" {
		return fmt.Errorf("error: user authentication required")
	}

	connStream.Close(websocket.StatusNormalClosure, "Authorized")

	return nil
}
