package models

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

// Pagination is a struct for pagination
type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	From        int64 `json:"from"`
	To          int64 `json:"to"`
}

type Payload struct {
	Items      interface{} `json:"items"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Response struct {
	Meta    Meta     `json:"meta"`
	Payload *Payload `json:"payload"`
}

func ResponseWrite(message string, data interface{}) *Response {
	return &Response{
		Meta: Meta{
			Code:    http.StatusOK,
			Message: message,
			Error:   false,
		},
		Payload: &Payload{
			Items: data,
		},
	}
}

func ResponseError(statusCode int, message string) *Response {
	return &Response{
		Meta: Meta{
			Code:    statusCode,
			Message: message,
			Error:   true,
		},
	}
}

func (response *Response) ToBytes() []byte {
	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// Paginate is a function for paginate
func (response *Response) Paginate(total, page, amount int64) (result *Response) {
	result = response
	response.Payload.Pagination = &Pagination{
		Total:       total,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(total) / float64(amount))),
		From:        (page-1)*amount + 1,
		To:          min(page*amount, total),
	}

	return
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
