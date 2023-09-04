package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting webserver on 8080")

	fs := http.FileServer(http.Dir("build"))
	http.Handle("/static", fs)

	panic(http.ListenAndServe(":8080", nil))
}
