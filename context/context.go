package context

import "os"

// Context represents the interface for querying information about the current environment
type Context interface {
	AuthToken() (string, error)
}

// New initializes a Context that reads from the filesystem
func New() Context {
	return &fsContext{}
}

// A Context implementation that queries the filesystem
type fsContext struct {
	authToken string
}

func (c *fsContext) AuthToken() (string, error) {
	if c.authToken != "" {
		return c.authToken, nil
	}

	token := os.Getenv("GITHUB_TOKEN")

	return token, nil
}
