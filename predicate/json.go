package predicate

import "encoding/json"

func JSON(value string) bool {
	return json.Valid([]byte(value))
}
