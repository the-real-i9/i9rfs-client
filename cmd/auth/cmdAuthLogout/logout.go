package cmdAuthLogout

import (
	"fmt"
	"i9rfs/client/globals"
)

func Execute() {
	globals.AppDataStore.RemoveItem("auth_jwt")
	globals.AppDataStore.RemoveItem("user")
	globals.AppDataStore.Save()

	fmt.Println("You've been logged out!")
}
