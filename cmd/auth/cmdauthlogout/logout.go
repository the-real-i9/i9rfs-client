package cmdauthlogout

import "i9pkgs/i9services"

func Execute() {
	i9services.LocalStorage.DeleteItem("auth_jwt")
}
