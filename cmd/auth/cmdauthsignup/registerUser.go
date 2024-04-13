package cmdauthsignup

import "nhooyr.io/websocket"

func registerUser(connStream *websocket.Conn) {
	for {
		// ask for username
		// ask for password

		// write user data along with signup_session_jwt to WS server

		// read response from connStream
		// if app_err, continue for loop and ask for the code again, else proceed to the next step and break

		// store user data and auth_jwt

		// close connStream
	}
}
