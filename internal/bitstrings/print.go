package bitstrings

import "fmt"

func RepresentPointerArray(in []*BitString128) string {
	res := ""
	for _, pointer := range in {
		res = fmt.Sprintf("%s | %d, %d", res, pointer.Upper(), pointer.Lower())
	}
	return fmt.Sprintf("[%s]", res)
}
