package cmdAuthLogout

import (
	"fmt"
	"i9rfs/client/appGlobals"
)

func Execute() {
	appGlobals.AppDataStore.RemoveItem("auth_jwt")
	appGlobals.AppDataStore.RemoveItem("user")
	appGlobals.AppDataStore.RemoveItem("i9rfs_work_path")
	appGlobals.AppDataStore.Save()

	fmt.Println("You've been logged out!")
}
