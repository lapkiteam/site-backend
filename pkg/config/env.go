package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func getValue(key string) (string, error) {
	file, err := os.Open(".env")

	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Split(scanner.Text(), "=")

		if readData[0] == key {
			return readData[1], nil
		}
	}

	return "", errors.New("key not found")
}

func GetAppUrl() (string, error) {
	appUrl, err := getValue("APP_URL")
	if err != nil {
		return "", err
	}

	return appUrl, nil
}
