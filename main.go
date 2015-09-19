package main

import (
	"log"
	"os"

	_ "github.com/dhiraj666/gocode/webfeed/matchers"
	"github.com/dhiraj666/gocode/webfeed/search"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("the")
}
