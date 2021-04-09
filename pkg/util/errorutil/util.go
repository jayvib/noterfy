package errorutil

import "errors"

// TryUnwrapErr tries unwrapping the root err if found
// otherwise it will return the passed error err.
func TryUnwrapErr(err error) error {
	if e := errors.Unwrap(err); e != nil {
		return e
	}
	return err
}
