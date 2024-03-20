package rest

type RegisterRequest struct {
	Date string
	Time string
}

type Response struct {
	Message string `json:"message"`
}
