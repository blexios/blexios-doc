package formatter

import (
	"fmt"
	"strings"
	"time"
)

func FormatDate(value any) string {
	if value == nil {
		return NullValue
	}

	switch v := value.(type) {
	case time.Time:
		return v.Format("02/01/2006")
	case string:
		s := strings.TrimSpace(v)
		for _, layout := range []string{"2006-01-02", "2006/01/02", "02/01/2006"} {
			if parsed, err := time.Parse(layout, s); err == nil {
				return parsed.Format("02/01/2006")
			}
		}
		return s
	default:
		return fmt.Sprint(v)
	}
}
