package vendor

import "fmt"

type OktaError struct {
	ErrorCode    string        `json:"errorCode"`
	ErrorSummary string        `json:"errorSummary"`
	ErrorLink    string        `json:"errorLink"`
	ErrorID      string        `json:"errorId"`
	ErrorCauses  []interface{} `json:"errorCauses"`
}

func (e *OktaError) ToString() string {
	return fmt.Sprintf("\nError Received From Okta:\nCode: [%v]\nSummary: [%v]\n\n", e.ErrorCode, e.ErrorSummary)
}

func (e *OktaError) Error() string {
	return fmt.Sprintf("\nError Received From Okta:\nCode: [%v]\nSummary: [%v]\n\n", e.ErrorCode, e.ErrorSummary)
}
