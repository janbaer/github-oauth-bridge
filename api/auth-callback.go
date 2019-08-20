package api

import (
	"fmt"
	"net/http"
)

// AuthCallback - Handles the authCallback request
func AuthCallback(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello from authCallback")

}
