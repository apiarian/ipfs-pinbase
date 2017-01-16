package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

type ErrorResponse struct {
	Message string `json:"error-message"`
	Details string `json:"error-details"`
}

func MarshalResponse(w http.ResponseWriter, c int, data interface{}) {
	j, err := json.Marshal(data)

	if err != nil {
		log.Printf("failed to marshal JSON: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)

		w.Write(
			[]byte(
				`{"error-message":"internal server error","error-details":"can't even."}`,
			),
		)

		return
	}

	w.WriteHeader(c)
	_, err = w.Write(j)

	if err != nil {
		log.Printf("failed to write response: %+v", err)
	}
}

func InternalServerError(w http.ResponseWriter) {
	MarshalResponse(
		w,
		http.StatusInternalServerError,
		ErrorResponse{
			Message: "internal server error",
			Details: "can't even.",
		},
	)
}
