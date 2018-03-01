package models

type Recipe struct {
	UId        int    `json:"uid" sql:"uid,pk"`
	Name       string `json:"name" sql:"name,notnull"`
	Vegetarian bool   `json:"vegetarian" sql:"vegetarian,notnull"`
	PrepTime   int    `json:"prep_time" sql:"prep_time,notnull"`
	Difficulty int    `json:"difficulty" sql:"difficulty,notnull"`
	UserId     int    `json:"user_id" sql:"user_id,notnull"`
	Ratings    []int  `json:"ratings" pg:"ratings,array"`
}

type Rating struct {
	Rating int `json:"rating"`
}
