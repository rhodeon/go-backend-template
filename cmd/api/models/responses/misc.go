package responses

type PingResponse struct {
	Body PingResponseBody
}

type PingResponseBody struct {
	Status string `json:"status" enum:"OK" doc:"Acknowledgement status"`
}
