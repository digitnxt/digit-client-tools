package build

import "fmt"

type StageError struct {
	Stage string
	Err   error
}

func (e *StageError) Error() string {
	return fmt.Sprintf("stage %s failed: %v", e.Stage, e.Err)
}

func (e *StageError) Unwrap() error {
	return e.Err
}
