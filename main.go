package main 

import (
	"fmt"
	"net/http"
	"power4/router"
)

func main() {
	r :=router.New()
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}