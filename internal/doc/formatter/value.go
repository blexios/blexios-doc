package formatter

import (
	"fmt"
	"strconv"
	"strings"
)

const NullValue = "—"

func FormatValue(value any, valueType string) string {
	if value == nil {
		return NullValue
	}

	switch strings.ToLower(valueType) {
	case "date":
		return FormatDate(value)
	case "time":
		return FormatTime(value)
	case "number":
		return FormatNumber(value)
	case "bool":
		return FormatBoolean(value)
	case "null":
		return NullValue
	default:
		return FormatString(value)
	}
}

func FormatString(value any) string {
	if value == nil {
		return NullValue
	}
	return fmt.Sprint(value)
}

func FormatNumber(value any) string {
	if value == nil {
		return NullValue
	}

	switch v := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	default:
		return fmt.Sprint(v)
	}
}

func FormatBoolean(value any) string {
	if value == nil {
		return NullValue
	}

	switch v := value.(type) {
	case bool:
		if v {
			return "Oui"
		}
		return "Non"
	default:
		return fmt.Sprint(v)
	}
}
