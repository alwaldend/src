package resources

import (
	"errors"
	"fmt"
	"strconv"
)

// resourceID accepts only canonical positive decimal resource identifiers.
func resourceID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse resource ID as a positive decimal integer: %w", err)
	}
	if id <= 0 || strconv.FormatInt(id, 10) != value {
		return 0, errors.New("resource ID must be a canonical positive decimal integer")
	}
	return id, nil
}
