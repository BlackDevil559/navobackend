package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/BlackDevil559/novahack2/routers"
)

func main() {
	r:=router.Router()
	fmt.Println("main.go file")
	log.Fatal(http.ListenAndServe(":4000",r))
	fmt.Println("Listening at 4000 Port")
}