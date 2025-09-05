package generator

import "fmt"

func Plural(one, few, many string) func(n int) string {
	return func(n int) string {
		mod10 := n % 10
		mod100 := n % 100

		out := many

		switch {
		case mod10 == 1 && mod100 != 11:
			out = one
		case mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14):
			out = few
		}

		return fmt.Sprintf("%d %s", n, out)
	}
}
