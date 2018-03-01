package models

type User struct {
	Id       int    `json:"uid" sql:"uid,pk"`
	Username string `json:"username" sql:"username,notnull"`
	IsActive bool   `json:"isactive" sql:"isactive,notnull"`
}

type Password struct {
	Hash string `json:"hash" sql:"hash,notnull"`
	Id   int    `json:"user_id" sql:"user_id,pk"`
}
