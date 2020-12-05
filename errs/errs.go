// Package errs implements error constants and also error wrapper of standard interface
package errs

import (
	"encoding/json"
	"errors"
)

var ErrHubNameIsempty = errors.New("hub name is empty")
var ErrHubNameOver100 = errors.New("hub name over 100 symbols")
var ErrHubIsAlreadyExists = errors.New("hub is already exists")
var ErrHubNameNotMatchedPattern = errors.New("the hub name should be contain only letters, digits, ., -, _")
var ErrHubNotFound = errors.New("hub not found")
var ErrNodeNameIsempty = errors.New("node name is empty")
var ErrNodeNameOver100 = errors.New("node name over 100 symbols")
var ErrNodeNotFound = errors.New("node not found")
var ErrNodeIsAlreadyExists = errors.New("node is already exists")
var ErrNodeNameNotMatchedPattern = errors.New("the node name should be contain only letters, digits, ., -, _")
var ErrPageNotFound = errors.New("page not found")

// Error wrapper of standard error interface. This entity include json marshaling
type Error struct {
	error string
}

// New Creates Error
func New(err error) *Error {
	return &Error{error: err.Error()}
}

func (e *Error) Error() string {
	return e.error
}

func (e *Error) MarshalJSON() ([]byte, error) {

	errormap := make(map[string]interface{})
	errormap["error"] = e.error

	return json.Marshal(errormap)

}
