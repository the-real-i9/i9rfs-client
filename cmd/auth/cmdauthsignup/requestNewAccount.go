package cmdauthsignup

import "nhooyr.io/websocket"

func requestNewAccount(connStream *websocket.Conn) {
	for {
		// ask for email

		// send email data to WS server

		// read response from connStream
		// store signup_session_jwt

		// if app_err, continue for loop and ask for the email again, else break
	}
}
