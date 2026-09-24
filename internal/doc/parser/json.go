package parser

import (
	"encoding/json"
	"errors"
)

var ErrEmptyArray = errors.New("the JSON array contains no records")
var ErrUnsupportedStructure = errors.New("unsupported JSON structure")

func ParseJSON(data []byte) ([]map[string]any, error) {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	items, ok := raw.([]any)
	if !ok {
		return nil, ErrUnsupportedStructure
	}

	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		row, ok := item.(map[string]any)
		if !ok {
			return nil, ErrUnsupportedStructure
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, ErrEmptyArray
	}

	return rows, nil
}
