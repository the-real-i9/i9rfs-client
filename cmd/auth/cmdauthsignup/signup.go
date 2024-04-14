package cmdauthsignup

import (
	"fmt"
	"i9pkgs/i9helpers"
	"i9rfs/client/cmd/rfsinter"
	"log"

	"nhooyr.io/websocket"
)

func Execute() {
	// connect to WS server
	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/auth/signup", "")
	if err != nil {
		log.Println(fmt.Errorf("signup error: %s", err))
	}

	defer connStream.CloseNow()

	signupSessionJwt, newAccEmail, err := requestNewAccount(connStream)
	if err != nil {
		return
	}

	if err := verifyEmail(connStream, signupSessionJwt, newAccEmail); err != nil {
		return
	}

	if err := registerUser(connStream, signupSessionJwt); err != nil {
		return
	}

	connStream.Close(websocket.StatusNormalClosure, "Signup success!")

	// signup is successful. At this point, the user is logged in
	// Therefore, we can Launch the Remote File System
	rfsinter.Launch()
}
