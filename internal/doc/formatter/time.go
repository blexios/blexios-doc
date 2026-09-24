package formatter

import (
	"fmt"
	"strings"
	"time"
)

func FormatTime(value any) string {
	if value == nil {
		return NullValue
	}

	switch v := value.(type) {
	case time.Time:
		return v.Format("15:04")
	case string:
		s := strings.TrimSpace(v)
		for _, layout := range []string{"15:04:05", "15:04", "15:04:05.999999999"} {
			if parsed, err := time.Parse(layout, s); err == nil {
				return parsed.Format("15:04")
			}
		}
		return s
	default:
		return fmt.Sprint(v)
	}
}
