package response

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Data any `json:"data,omitempty"`
}

// JSON writes json response
func JSON(w http.ResponseWriter, code int, body any) {
	// check body
	if body == nil {
		w.WriteHeader(code)
		return
	}

	response := apiResponse{
		Data: body,
	}

	// marshal body
	bytes, err := json.Marshal(response)
	if err != nil {
		// default error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set header
	w.Header().Set("Content-Type", "application/json")

	// set status code
	w.WriteHeader(code)

	// write body
	w.Write(bytes)
}
