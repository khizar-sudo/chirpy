package main

import (
	"github.com/khizar-sudo/chirpy/handlers"
	_ "github.com/lib/pq"
)

func main() {
	handlers.Init()
}
