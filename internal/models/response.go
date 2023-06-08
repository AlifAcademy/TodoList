package models

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// Meta is a struct information response
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

// Payload is a struct for Date Payload
type Payload struct {
	Items      interface{} `json:"items"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Response is a struct for response
type Response struct {
	Meta    Meta     `json:"meta"`
	Payload *Payload `json:"payload"`
}

// ResponseWrite ....
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

// ResponseError ....
func ResponseError(statusCode int, message string) *Response {
	return &Response{
		Meta: Meta{
			Code:    statusCode,
			Message: message,
			Error:   true,
		},
	}
}

// ToBytes convert to []byte message Response
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
