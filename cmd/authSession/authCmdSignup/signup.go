package authCmdSignup

import (
	"fmt"
	"i9rfs/client/cmd/mgmtSession"
	"i9rfs/client/helpers"
	"i9rfs/client/services/authServices"
	"log"

	"nhooyr.io/websocket"
)

func Execute() {
	// connect to WS server
	connStream, err := helpers.WSConnect("ws://localhost:8000/api/auth/signup", "")
	if err != nil {
		log.Printf("signup: wsconn error: %s\n", err)
		return
	}

	defer connStream.CloseNow()

	signupSessionJwt, newAccEmail, err := authServices.RequestNewAccount(connStream)
	if err != nil {
		fmt.Println(err)
		return
	}

	signupSessionJwt2, err := authServices.VerifyEmail(connStream, signupSessionJwt, newAccEmail)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := authServices.RegisterUser(connStream, signupSessionJwt2); err != nil {
		fmt.Println(err)
		return
	}

	connStream.Close(websocket.StatusNormalClosure, "Signup success!")

	// signup is successful. At this point, the user is logged in
	// Therefore, we can Begin the Remote File System session
	mgmtSession.Begin()
}
