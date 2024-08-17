package responses

type PingResponse struct {
	Body struct {
		Status string `json:"status" enum:"OK" doc:"Acknowledgement status"`
	}
}
