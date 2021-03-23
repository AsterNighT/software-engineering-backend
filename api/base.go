package api

type ReturnedData struct {
	// A simple string indicating the status.
	// Is it ok, or some error occurs? If so, what is the error?
	// It should be "ok" is everything goes fine
	Status string `json:"status" `

	// Anything you want to pass to the frontend, but make it simple and necessary
	// If there's nothing to return, this field will be omitted
	Data interface{} `json:"data,omitempty"`
}

func Return(status string, data interface{}) ReturnedData {
	return ReturnedData{Status: status, Data: data}
}
