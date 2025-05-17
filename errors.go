package mgmgo

import "fmt"

type ErrMACsDiffer struct {
	err       string
	Computed  []byte
	Presented []byte
}

func (e *ErrMACsDiffer) Error() string {
	return fmt.Sprintf("MACs are different: \nComputed:  %v, \nPresented: %v, \nError: %s", e.Computed, e.Presented, e.err)
}
