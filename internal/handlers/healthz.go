package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	// but we want json apis right
	w.Header().Set("Content-Type", "application/json")

	// w.Write([]byte("all ok")) // by default we are sending text/plain content-type

	w.WriteHeader(http.StatusOK) // this line should comes after all Header set because it sends the header already to the respjnse

	w.Write([]byte(`{"status" : "okay"}`))
}
