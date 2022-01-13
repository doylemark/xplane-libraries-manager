package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

const (
	cookiesFile = "cookies.txt"
	userFile    = "user.txt"
)

func storeUser(user User) error {
	f, err := os.OpenFile(userFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	out := user.name + " " + user.password

	_, err = f.WriteString(out)

	if err != nil {
		return err
	}

	return nil
}

func readUser() (User, error) {
	var user User

	b, err := os.ReadFile(userFile)

	if err != nil {
		return user, err
	}

	parts := strings.Split(string(b), " ")

	user.name = parts[0]
	user.password = parts[1]

	return user, nil
}

func storeCookies(cookies []*http.Cookie) error {
	f, err := os.OpenFile(cookiesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	for _, cookie := range cookies {
		_, err = f.WriteString(cookie.Name + " " + cookie.Value + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func readCookies() []*http.Cookie {
	var cookies []*http.Cookie

	f, err := os.Open("cookies.txt")

	if err != nil {
		return cookies
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		cookies = append(cookies, &http.Cookie{Name: parts[0], Value: parts[1]})
	}

	return cookies
}
