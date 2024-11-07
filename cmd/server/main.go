package main

import (
	api "cast/internal/api/router"
)

func main() {
	r := api.InitRouter()
	r.Run()
}
