package main

import (
	"forum/logic"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	logic.InitData()
	logic.LaunchApp()
}
