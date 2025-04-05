package main

import "github.com/dimfu/mbinder/cmd"

func main() {
	cmd, err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	db, err := cmd.DB.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
