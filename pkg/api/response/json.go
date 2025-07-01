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

	bytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	w.Write(bytes)
}
