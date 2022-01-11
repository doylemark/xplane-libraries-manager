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
	"net/http"

	"github.com/chromedp/cdproto/network"
	dp "github.com/chromedp/chromedp"
)

const (
	loginUrl        = "https://forums.x-plane.org/index.php?/login"
	usernameInputId = "auth"
	passwordInputId = "password"
	signInBtnId     = "elSignIn_submit"
)

func Login(username string, password string) []*http.Cookie {
	var loginCookies []*http.Cookie

	ctx, cancel := dp.NewContext(context.Background(), dp.WithBrowserOption())
	defer cancel()

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
				return err
			}

			for _, cookie := range cookies {
				loginCookies = append(loginCookies, &http.Cookie{
					Name:  cookie.Name,
					Value: cookie.Value,
				})
			}

			return nil
		}),
	)

	return loginCookies
}

func MakeAuthorizedGet(cookies []*http.Cookie, url string) (*http.Response, error) {
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
