package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
