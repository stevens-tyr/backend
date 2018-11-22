package main

import (
	"backend/api"
)

func main() {
	server := api.SetUp()
	server.Run(":5555")
}
