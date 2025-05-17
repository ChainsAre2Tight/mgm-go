package bitstrings

import "fmt"

func RepresentPointerArray(in []*BitString128) string {
	res := ""
	for _, pointer := range in {
		res = fmt.Sprintf("%s | %0.16x, %0.16x", res, pointer.Upper(), pointer.Lower())
	}
	return fmt.Sprintf("[%s]", res)
}
