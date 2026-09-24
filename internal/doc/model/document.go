package model

import (
	"time"
)

type Metadata struct {
	SourcePath    string    `json:"source_path"`
	GeneratedAt   time.Time `json:"generated_at"`
	RecordCount   int       `json:"record_count"`
	ColumnCount   int       `json:"column_count"`
	Orientation   string    `json:"orientation"`
	DocumentTitle string    `json:"document_title"`
}

type Document struct {
	Title    string   `json:"title"`
	Columns  []Column `json:"columns"`
	Rows     []Row    `json:"rows"`
	Metadata Metadata `json:"metadata"`
}

func NewDocument(title string, columns []Column, rows []map[string]any, sourcePath string, orientation string) Document {
	document := Document{
		Title:    title,
		Columns:  columns,
		Rows:     make([]Row, 0, len(rows)),
		Metadata: Metadata{SourcePath: sourcePath, GeneratedAt: time.Now(), RecordCount: len(rows), ColumnCount: len(columns), Orientation: orientation, DocumentTitle: title},
	}

	for _, row := range rows {
		document.Rows = append(document.Rows, Row{Values: cloneMap(row)})
	}

	return document
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}

	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
