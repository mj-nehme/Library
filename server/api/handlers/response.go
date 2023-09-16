package handlers

import "encoding/json"

type ErrorResponse struct {
	Error string `json:"error"`
}

func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"error": e.Error,
	})
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (m MessageResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": m.Message,
	})
}

type StatusResponse struct {
	Status string `json:"status"`
}

func (s StatusResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"status": s.Status,
	})
}
