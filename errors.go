package migadu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Error struct {
	StatusCode int
	Code       string `json:"error"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("migadu: %s: %s", e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("migadu: unexpected status %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("migadu: unexpected status %d", e.StatusCode)
}

func newError(resp *http.Response) error {
	e := &Error{StatusCode: resp.StatusCode}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, e); err != nil || e.Code == "" {
		e.Message = strings.TrimSpace(string(body))
	}
	return e
}
