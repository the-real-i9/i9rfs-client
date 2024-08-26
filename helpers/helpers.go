package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

func WSConnect(url string, authToken string) (connStream *websocket.Conn, err error) {
	dialOptions := &websocket.DialOptions{
		HTTPHeader: make(http.Header),
	}

	if authToken != "" {
		dialOptions.HTTPHeader.Set("Authorization", authToken)
	}

	conn, _, err := websocket.Dial(context.Background(), url, dialOptions)

	if err != nil {
		return nil, fmt.Errorf("WSConnect: websocket.Dial: %s", err)
	}

	return conn, nil
}

func ParseTo(input any, output any) {
	bt, _ := json.Marshal(input)

	json.Unmarshal(bt, output)
}

func ToJSON(val any) string {
	bt, _ := json.Marshal(val)

	return string(bt)
}
