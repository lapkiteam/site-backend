package auth

import (
	"github.com/goccy/go-json"
	"os"
)

type User struct {
	Data []struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"data"`
}

func CheckLogin(login, password string) bool {
	fileData, err := os.ReadFile("auth.json")
	if err != nil {
		return false
	}

	var user User

	err = json.Unmarshal(fileData, &user)
	if err != nil {
		return false
	}

	for _, user := range user.Data {
		if login == user.Login && user.Password == password {
			return true
		}
	}

	return false
}
