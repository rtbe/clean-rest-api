package handlers

// package handlers REST API
//
// Documentation of web API
//
//  Schemes: http
//  Host: localhost
//  Base path: /
//  Version: 1.0.0
//
//  Consumes:
//  	- application/json
//
//  Produces:
//  	- application/json
//
//  swagger:meta

// Error response
// swagger:response errorResponse
type errorResponse struct {
	// Error response message
	//
	// in: body
	Body struct {
		// Example: There are some error
		//
		Error  string `json:"error"`
		Fields string `json:"fields,omitempty"`
	}
}

// Empty response
// swagger:response errorResponse
type emptyResponse struct {
}

// Status response
type statusResponse struct {
	// Statusresponse message
	//
	// in: body
	Body struct {
		// Example: up
		//
		Status string `json:"status"`
		Host   string `json:"host"`
	}
}
