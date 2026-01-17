package tc

import (
	"time"
	"errors"
	"strconv"
)

type Textcat struct {
	Function LogicBridge
	Sessions *SessionManager
}

type Session struct {
	Username string
	Connection any
}

var ErrNotFound = errors.New("data does not found")

func (tc *Textcat) CreateUser(username string, password string) error {
	tc.Function.LogMsg("info", "user register", username, password)

	tblErr := tc.Function.CreateTable("users")
	if tblErr != nil {
		return MakeError("server_error:", tblErr)
	}

	// Validate username
	if !IsValidUsername(username) {
		return errors.New("error: invalid username: must be alphanumeric, can contain '-' and '_'")
	}

	// Build user object
	user := User{
		Username:  username,
		Password:  password,
		Created:   time.Now(),
		LastLogin: time.Now(),
	}

	// Delegate existence check to the Handler
	var existing User
	err := tc.Function.GetData("users", func(v any) bool {
		u := v.(*User)
		return u.Username == username // return the sub-func() if user exists
	}, &existing)
	

	if err == nil {
		return errors.New("error: username is taken")
	}

	// Delegate storage to the Handler
	err = tc.Function.StoreData("users", user)
	if err != nil {
		return MakeError("server_error", err)
	}
	return errors.New("ok")
}

func (tc *Textcat) LoginUser(username string, password string, connection RequestWriter) error {
	tc.Function.LogMsg("info", "user login", username, password)

	// Validate username
	if !IsValidUsername(username) {
		return errors.New("error: invalid username: must be alphanumeric, can contain '-' and '_'")
	}

	var existing User
	err := tc.Function.GetData("users", func(v any) bool {
		u := v.(*User)
		return u.Username == username
	}, &existing)

	if errors.Is(err, ErrNotFound) {
		return errors.New("error: user does not exist")
	}
	if err != nil {
		return MakeError("server_error:", err)
	}

	if existing.Password != password {
		return errors.New("error: invalid password")
	}

	// After login succeeds
	existing.LastLogin = time.Now()

	// Store updated user
	if err := tc.Function.StoreData("users", existing); err != nil {
		return MakeError("server_error:", "failed to update last login:", err)
	}

	UserSession := Session{
		Username: username,
		Connection: nil,
	}

	userToken := tc.Sessions.GetUnused()

	tc.Sessions.Add(userToken, UserSession)

	tc.Function.MakeRequest("status", "YourToken", strconv.Itoa(userToken), "ok", connection)
	return errors.New("ok")
}