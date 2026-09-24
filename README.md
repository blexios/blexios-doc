# Blexios DOC

Blexios DOC is a Go CLI for converting raw JSON data into a professional PDF document.

## Usage

```powershell
./blexios.exe doc:generate --pdf --path=examples/doc/presences.json --title="Rapport de présence"
```

## Technical flow

The command is handled in the CLI entrypoint and then routed through the document generation pipeline:

```text
CLI arguments
  ↓
JSON file validation
  ↓
Array-of-objects detection
  ↓
Column analysis and label generation
  ↓
Internal document model
  ↓
Value formatting
  ↓
PDF renderer
  ↓
Generated file
```

## Main components

- `main.go` : entry point for the binary
- `cmd/doc/generate.go` : command parsing and flags
- `internal/doc/service/generator.go` : orchestration layer
- `internal/doc/parser/` : JSON parsing and structure detection
- `internal/doc/model/` : internal document model
- `internal/doc/formatter/` : date, time, null, string formatting
- `internal/doc/renderer/pdf/` : PDF rendering

## Example

```powershell
./blexios.exe doc:generate --pdf --path=examples/doc/products.json --output=examples/doc/custom_report.pdf --title="Catalogue"
```

## Features

- Detects JSON arrays of objects
- Auto-detects columns and labels
- Handles null values as `—`
- Formats dates and times
- Generates PDF automatically
- Keeps the internal model separate from the PDF renderer

## Development

```bash
go test ./...
go build -o blexios.exe .
```

## More documentation

- [docs/doc-generator.md](docs/doc-generator.md)
- [docs/technical-doc.md](docs/technical-doc.md)
