package cmdauthlogout

import (
	"fmt"
	"i9pkgs/i9services"
)

func Execute() {
	i9services.LocalStorage.DeleteItem("auth_jwt", "user")

	fmt.Println("You've been logged out!")
}
