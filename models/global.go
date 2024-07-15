package models

import "time"

type Response struct {
	StatusCode       string      `json:"statusCode"`
	Success          bool        `json:"success"`
	ResponseDatetime time.Time   `json:"responseDatetime"`
	Result           interface{} `json:"result"`
	Message          string      `json:"message"`
}

type ResponseList struct {
	StatusCode       string      `json:"statusCode"`
	Success          bool        `json:"success"`
	ResponseDatetime time.Time   `json:"responseDatetime"`
	CountData        int         `json:"countData"`
	Result           interface{} `json:"result"`
	Message          string      `json:"message"`
}

type ErrorMsg struct {
	Status    string `json:"status"`
	ErrorDesc string `json:"errordesc"`
}

type RequestList struct {
	Order   string `json:"order" validate:"required"` //asc atau desc
	OrderBy string `json:"orderBy"`                   //order by per field
	Limit   int    `json:"limit" validate:"required"` //jumlah per page
	Page    int    `json:"page" validate:"required"`  //Halaman ke brp yg mau di load
	Keyword string `json:"keyword"`                   //untuk search
}

type PageOffsetLimit struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
