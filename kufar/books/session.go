package main

import (
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 5,
		HttpOnly: true,
	}
	os.Setenv("SESSION_AUTH_KEY", string(authKeyOne))
	os.Setenv("SESSION_ENCRYPT_KEY", string(encryptionKeyOne))
}
