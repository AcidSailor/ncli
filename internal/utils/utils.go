package utils

import (
	"fmt"
	"strings"
)

func NetconfStrip(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "]]>]]>")
	return s
}

func FlatPathToSubtreeWithValue(s string, v string) string {
	sa := strings.FieldsFunc(s, func(r rune) bool {
		return r == '/'
	})

	f := make([]string, len(sa)*2+1)

	for i, elem := range sa {
		f[i] = fmt.Sprintf("<%s>", elem)
		if i == len(sa)-1 {
			f[i+1] = v
		}
		f[len(f)-1-i] = fmt.Sprintf("</%s>", elem)
	}

	return strings.Join(f, "")
}

func WrapWithTags(s string, tag string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, s, tag)
}
