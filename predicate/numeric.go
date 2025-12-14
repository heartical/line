package predicate

import "strconv"

func Integer(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func Number(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
