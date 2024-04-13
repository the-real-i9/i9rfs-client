package cmdauthsignup

import (
	"fmt"
	"i9pkgs/i9helpers"
	"log"
)

func Execute() {
	// connect to WS server
	connStream, err := i9helpers.WSConnect("ws://localhost:8000/api/auth/signup", "")
	if err != nil {
		log.Println(fmt.Errorf("signup error: %s", err))
	}

	defer connStream.CloseNow()

	requestNewAccount(connStream)

	verifyEmail(connStream)

	registerUser(connStream)
}
