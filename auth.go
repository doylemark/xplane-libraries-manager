/*
 * CDDL HEADER START
 *
 * This file and its contents are supplied under the terms of the
 * Common Development and Distribution License ("CDDL"), version 1.0.
 * You may only use this file in accordance with the terms of version
 * 1.0 of the CDDL.
 *
 * A full copy of the text of the CDDL should have accompanied this
 * source.  A copy of the CDDL is also available via the Internet at
 * http://www.illumos.org/license/CDDL.
 *
 * Copyright 2022 Mark Doyle. All rights reserved.
 */

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/cdproto/network"
	dp "github.com/chromedp/chromedp"
)

const (
	loginUrl        = "https://forums.x-plane.org/index.php?/login"
	loginFailedMsg  = "The display name, email address or password was incorrect"
	usernameInputId = "auth"
	passwordInputId = "password"
	signInBtnId     = "elSignIn_submit"
	cookiesFile     = "cookies.txt"
	userFile        = "user.txt"
)

func login(username, password string) ([]*http.Cookie, error) {
	var loginCookies []*http.Cookie

	opts := append(dp.DefaultExecAllocatorOptions[:],
		dp.Flag("headless", true),
		dp.Flag("disable-gpu", false),
		dp.Flag("enable-automation", false),
		dp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := dp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := dp.NewContext(allocCtx)
	defer cancel()
	fmt.Println(username, password)
	dp.Run(ctx,
		dp.Navigate(loginUrl),
		dp.WaitReady("body"),

		dp.SetValue(usernameInputId, username, dp.ByID),
		dp.SetValue(passwordInputId, password, dp.ByID),
		dp.Click(signInBtnId, dp.ByID),

		dp.WaitReady("body"),
		dp.ActionFunc(func(ctx context.Context) error {

			cookies, err := network.GetAllCookies().Do(ctx)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			for _, cookie := range cookies {
				fmt.Println(cookie.Name)
				loginCookies = append(loginCookies, &http.Cookie{
					Name:  cookie.Name,
					Value: cookie.Value,
				})
			}

			return nil
		}),
	)

	if len(loginCookies) <= 5 {
		return loginCookies, errors.New("login failed, not enough cookies")
	}

	storeCookies(loginCookies)
	return loginCookies, nil
}

func makeAuthorizedGet(cookies []*http.Cookie, url string) (*http.Response, error) {
	var client http.Client
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func storeUser(username, password string) error {
	f, err := os.OpenFile(userFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	out := username + " " + password

	_, err = f.WriteString(out)

	if err != nil {
		return err
	}

	return nil
}

func readUser() (*string, *string) {
	b, err := os.ReadFile(userFile)

	if err != nil {
		return nil, nil
	}

	parts := strings.Split(string(b), " ")

	return &parts[0], &parts[1]
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
