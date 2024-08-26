package initializers

import (
	"fmt"
	"i9rfs/client/globals"
	"os"
	"os/exec"
)

func initAppDataStore() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	appDir := fmt.Sprintf("%s/.i9rfs", homeDir)

	exec.Command("mkdir", appDir).Run()

	return globals.AppDataStore.Revive(fmt.Sprintf("%s/localStorage.json", appDir))
}

func InitApp() error {
	if err := initAppDataStore(); err != nil {
		return err
	}

	return nil
}
