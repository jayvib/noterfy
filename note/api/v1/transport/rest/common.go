package rest

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"noterfy/note"
	"noterfy/pkg/util/errorutil"
	"strconv"
)

// StatusClientClosed is an http status where the client cancels a request.
const StatusClientClosed = 499

func newErrorWrapper(err error) errorWrapper {
	return errorWrapper{
		origErr:    err,
		message:    getMessage(err),
		statusCode: getStatusCode(err),
	}
}

type errorWrapper struct {
	origErr    error
	message    string
	statusCode int
}

func (e errorWrapper) error() error {
	return errorutil.TryUnwrapErr(e.origErr)
}

func (e errorWrapper) Error() string {
	return e.origErr.Error()
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorWrapper)
	if ok && e.error() != nil {
		encodeError(e, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ew errorWrapper, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(ew.statusCode)

	logrus.Error(ew.origErr)

	_ = json.NewEncoder(w).Encode(ResponseError{
		Message: ew.message,
	})
}

func getStatusCode(err error) (statusCode int) {
	err = errorutil.TryUnwrapErr(err)
	switch err {
	case note.ErrNotFound:
		statusCode = http.StatusNotFound
	case note.ErrNilID:
		statusCode = http.StatusBadRequest
	case note.ErrExists:
		statusCode = http.StatusConflict
	case note.ErrCancelled:
		statusCode = StatusClientClosed
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}

func getMessage(err error) (message string) {
	causeErr := errorutil.TryUnwrapErr(err)
	switch causeErr {
	case note.ErrExists:
		message = "Note already exists"
	case note.ErrCancelled, context.Canceled:
		message = "Request cancelled"
	case note.ErrNotFound:
		message = "Note not found"
	case note.ErrNilID:
		message = "Empty note identifier"
	default:
		message = "Unexpected error"
	}
	return
}

func convertAtoU(s string) uint64 {
	v, _ := strconv.Atoi(s)
	return uint64(v)
}
