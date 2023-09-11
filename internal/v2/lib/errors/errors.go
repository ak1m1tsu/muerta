package errors

import "fmt"

func WithOp(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
