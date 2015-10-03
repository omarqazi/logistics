package main

import (
	"fmt"
	"net/http"
)

type LogisticsRouter struct {
}

func (l LogisticsRouter) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintln(w,"hello world")
}

