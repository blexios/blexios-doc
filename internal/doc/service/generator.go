package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"blexios/internal/doc/model"
	"blexios/internal/doc/parser"
	pdfRenderer "blexios/internal/doc/renderer/pdf"
)

type Options struct {
	PDF         bool
	Path        string
	Output      string
	Title       string
	Template    string
	Orientation string
	Verbose     bool
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(opts Options) (*model.Document, string, error) {
	if opts.Path == "" {
		return nil, "", errors.New("input file not found:\n" + opts.Path)
	}

	info, err := os.Stat(opts.Path)
	if err != nil || info.IsDir() {
		return nil, "", fmt.Errorf("input file not found:\n%s", opts.Path)
	}

	raw, err := os.ReadFile(opts.Path)
	if err != nil {
		return nil, "", fmt.Errorf("unable to read input file:\n%s", opts.Path)
	}

	var root any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, "", fmt.Errorf("invalid JSON.\nUnable to parse:\n%s", opts.Path)
	}

	rows, err := parser.ParseJSON(raw)
	if err != nil {
		if errors.Is(err, parser.ErrEmptyArray) {
			return nil, "", errors.New("the JSON array contains no records.")
		}
		if errors.Is(err, parser.ErrUnsupportedStructure) {
			return nil, "", errors.New("unsupported JSON structure.\nExpected an array of objects.")
		}
		return nil, "", fmt.Errorf("invalid JSON.\nUnable to parse:\n%s", opts.Path)
	}

	columns, err := parser.DetectColumns(rows)
	if err != nil {
		return nil, "", errors.New("unsupported JSON structure.\nExpected an array of objects.")
	}

	if opts.Title == "" {
		opts.Title = buildDefaultTitle(opts.Path)
	}

	orientation := opts.Orientation
	if strings.EqualFold(orientation, "auto") || orientation == "" {
		if len(columns) > 5 {
			orientation = "landscape"
		} else {
			orientation = "portrait"
		}
	}

	doc := model.NewDocument(opts.Title, columns, rows, opts.Path, orientation)
	doc.Metadata.Orientation = orientation

	outputPath, err := resolveOutputPath(opts.Path, opts.Output)
	if err != nil {
		return nil, "", err
	}

	renderer := pdfRenderer.NewRenderer(orientation)
	if err := renderer.Render(doc, outputPath); err != nil {
		return nil, "", errors.New("unable to generate PDF.")
	}

	return &doc, outputPath, nil
}

func buildDefaultTitle(inputPath string) string {
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	base = strings.TrimSuffix(base, ext)
	if base == "" {
		return "Blexios Document"
	}
	return titleize(base)
}

func titleize(value string) string {
	if value == "" {
		return "Blexios Document"
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '-' || r == '_' || r == '.' || r == '/' || r == ' '
	})
	if len(parts) == 0 {
		return "Blexios Document"
	}

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.EqualFold(part, "id") {
			result = append(result, "ID")
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result = append(result, string(runes))
	}
	return strings.Join(result, " ")
}

func resolveOutputPath(inputPath, output string) (string, error) {
	if strings.TrimSpace(output) != "" {
		return uniqueOutputPath(output), nil
	}

	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	candidate := base + ".pdf"
	return uniqueOutputPath(candidate), nil
}

func uniqueOutputPath(target string) string {
	if _, err := os.Stat(target); err != nil {
		return target
	}

	ext := filepath.Ext(target)
	base := strings.TrimSuffix(target, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(candidate); err != nil {
			return candidate
		}
	}
}
