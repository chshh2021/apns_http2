package apns_http2

import (
	"encoding/json"
	"io"
	"time"
)

func parseErrorResponse(body io.Reader, statusCode int) error {
	var response struct {
		Reason    string `json:"reason"`
		Timestamp int64  `json:"timestamp"`
	}
	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		return err
	}

	es := &Error{
		Reason: mapErrorReason(response.Reason),
		Status: statusCode,
	}

	if response.Timestamp != 0 {
		es.Timestamp = time.Unix(response.Timestamp/1000, 0)
	}
	return es
}
