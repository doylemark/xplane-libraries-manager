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
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
)

type User struct {
	name     string
	password string
}

func refreshAuth() []*http.Cookie {
	cookies := readCookies()

	cookiesValid := checkCookies(cookies)

	if cookiesValid {
		return cookies
	}

	user, err := readUser()

	if err != nil {
		return nil
	}

	newCookies, err := refreshCookies(user)
	cookiesValid = checkCookies(newCookies)

	if !cookiesValid || err != nil {
		return newCookies
	}

	return nil
}

func refreshCookies(user User) ([]*http.Cookie, error) {
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

	dp.Run(ctx,
		dp.Navigate(loginUrl),
		dp.WaitReady("body"),

		dp.SetValue(usernameInputId, user.name, dp.ByID),
		dp.SetValue(passwordInputId, user.password, dp.ByID),
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

func checkCookies(cookies []*http.Cookie) bool {
	resp, err := makeAuthorizedGet(cookies, "https://forums.x-plane.org/")

	if err != nil {
		fmt.Println(err)
		return false
	}

	d, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return false
	}

	str := string(d)

	return !strings.Contains(str, "Existing user? Sign In")

}
