package main
import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting logistics server...")
	if err := http.ListenAndServe(":8080",nil); err != nil {
		fmt.Println("Error starting logistics server:",err)
		return
	}
}