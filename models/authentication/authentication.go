package models

import "time"

type Authentication struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Jwt struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"exception"`
}

type TokenClaims struct {
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
	UserId   int       `json:"userid"`
}
