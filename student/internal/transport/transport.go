package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type (
	CreateUpdateStudentRequest struct {
		ID          string
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		StudentRoom int    `json:"room"`
	}
	CreateUpdateStudentResponse struct {
		Ok string `json:"ok"`
	}

	GetStudentRequest struct {
		ID string `json:"id"`
	}
	GetStudentResponse struct {
		ID          string `json:"id"`
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		StudentRoom int    `json:"room"`
	}

	DeleteStudentResponse struct {
		Ok string `json:"ok"`
	}
)

type URLParams struct {
	ID string
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateUpdateStudentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUpdateStudentRequest

	id := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ID = id

	return req, nil
}

func decodeGetDeleteStudentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetStudentRequest
	id := chi.URLParam(r, "id")

	req = GetStudentRequest{
		ID: id,
	}

	return req, nil
}
