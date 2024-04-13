package cmdauthsignup

import "nhooyr.io/websocket"

func verifyEmail(connStream *websocket.Conn) {
	for {
		// ask for verf code

		// send code along with signup_session_jwt to WS server

		// read response from connStream

		// if app_err, continue for loop and ask for the code again, else break
	}
}
