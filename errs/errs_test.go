package errs

import (
	"encoding/json"
	"testing"
)

func TestErrs(t *testing.T) {

	e := New(ErrNodeNotFound)

	if ErrNodeNotFound.Error() != e.Error() {
		t.Errorf("error must be %s", ErrNodeNotFound.Error())
	}

	_, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}

}
