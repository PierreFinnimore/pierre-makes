package types

import "github.com/a-h/templ"

// AuthUser represents an user that might be authenticated.
type AuthUser struct {
	ID       int
	Email    string
	LoggedIn bool
}

// Check should return true if the user is authenticated.
// See handlers/auth.go.
func (user AuthUser) Check() bool {
	return user.ID > 0 && user.LoggedIn
}

type Nav struct {
	Label    string
	Selected bool
	Link     templ.SafeURL
}
