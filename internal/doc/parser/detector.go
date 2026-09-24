package parser

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"blexios/internal/doc/model"
)

func DetectColumns(rows []map[string]any) ([]model.Column, error) {
	if len(rows) == 0 {
		return nil, ErrEmptyArray
	}

	order := make([]string, 0)
	seen := map[string]bool{}
	for _, row := range rows {
		for key := range row {
			if seen[key] {
				continue
			}
			seen[key] = true
			order = append(order, key)
		}
	}

	sort.Slice(order, func(i, j int) bool {
		left, right := order[i], order[j]
		leftKey, rightKey := columnPriority(left), columnPriority(right)
		if leftKey != rightKey {
			return leftKey < rightKey
		}
		return strings.ToLower(left) < strings.ToLower(right)
	})

	columns := make([]model.Column, 0, len(order))
	for _, key := range order {
		values := make([]any, 0, len(rows))
		for _, row := range rows {
			if value, exists := row[key]; exists {
				values = append(values, value)
			}
		}

		valueType := InferValueType(values)
		columns = append(columns, model.Column{
			Key:   key,
			Label: HumanizeKey(key),
			Type:  valueType,
		})
	}

	return columns, nil
}

func columnPriority(key string) int {
	lowered := strings.ToLower(key)
	if strings.Contains(lowered, "id") {
		return 0
	}
	if strings.Contains(lowered, "date") {
		return 1
	}
	if strings.Contains(lowered, "name") {
		return 2
	}
	return 3
}

func InferValueType(values []any) string {
	if len(values) == 0 {
		return "string"
	}

	allNil := true
	allBool := true
	allNumber := true
	allDate := true
	allTime := true

	for _, value := range values {
		if value == nil {
			continue
		}
		allNil = false

		switch value.(type) {
		case bool:
			allNumber = false
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			allBool = false
		case string:
			s := value.(string)
			allBool = false
			allNumber = false
			if !isDateString(s) {
				allDate = false
			}
			if !isTimeString(s) {
				allTime = false
			}
		default:
			allBool = false
			allNumber = false
			allDate = false
			allTime = false
		}
	}

	if allNil {
		return "null"
	}
	if allBool {
		return "bool"
	}
	if allNumber {
		return "number"
	}
	if allDate {
		return "date"
	}
	if allTime {
		return "time"
	}

	return "string"
}

func isDateString(value string) bool {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return false
	}

	if _, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return false
	}

	layouts := []string{"2006-01-02", "2006/01/02", "02/01/2006"}
	for _, layout := range layouts {
		if _, err := time.Parse(layout, cleaned); err == nil {
			return true
		}
	}

	return false
}

func isTimeString(value string) bool {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return false
	}

	layouts := []string{"15:04", "15:04:05", "15:04:05.999999999"}
	for _, layout := range layouts {
		if _, err := time.Parse(layout, cleaned); err == nil {
			return true
		}
	}
	return false
}

func HumanizeKey(key string) string {
	cleaned := strings.TrimSpace(key)
	if cleaned == "" {
		return ""
	}

	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	cleaned = strings.ReplaceAll(cleaned, ".", " ")
	cleaned = strings.ReplaceAll(cleaned, "/", " ")

	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	cleaned = re.ReplaceAllString(cleaned, `${1} ${2}`)

	parts := strings.Fields(cleaned)
	if len(parts) == 0 {
		return ""
	}

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		word := strings.TrimSpace(part)
		if word == "" {
			continue
		}

		if strings.EqualFold(word, "id") {
			result = append(result, "ID")
			continue
		}
		if strings.EqualFold(word, "url") {
			result = append(result, "URL")
			continue
		}

		r := []rune(word)
		if len(r) == 0 {
			continue
		}
		formatted := string(unicode.ToUpper(r[0])) + string(r[1:])
		result = append(result, formatted)
	}

	return strings.Join(result, " ")
}
