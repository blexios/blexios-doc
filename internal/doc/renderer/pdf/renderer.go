package pdf

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"blexios/internal/doc/formatter"
	"blexios/internal/doc/model"
)

type Renderer struct {
	Orientation string
}

func NewRenderer(orientation string) *Renderer {
	if orientation == "" {
		orientation = "auto"
	}
	return &Renderer{Orientation: orientation}
}

func (r *Renderer) Render(document model.Document, outputPath string) error {
	orientation := determineOrientation(document.Columns, r.Orientation)
	pdf := gofpdf.New(orientation, "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.SetAutoPageBreak(true, 12)

	fontName := "Helvetica"
	if fontPath := resolveFontPath(); fontPath != "" {
		pdf.AddUTF8Font("DejaVuSans", "", fontPath)
		fontName = "DejaVuSans"
	}
	pdf.SetFont(fontName, "", 11)
	pdf.AddPage()

	pdf.SetTextColor(25, 25, 25)
	pdf.SetFont(fontName, "", 18)
	pdf.MultiCell(0, 9, document.Title, "", "L", false)
	pdf.Ln(3)

	pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 6, fmt.Sprintf("Source: %s", document.Metadata.SourcePath), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Date: %s", document.Metadata.GeneratedAt.Format("02/01/2006 15:04")), "", 1, "L", false, 0, "")
	pdf.Ln(4)

	widths := columnWidths(pdf, len(document.Columns))
	drawHeader(pdf, document.Columns, widths, fontName)
	drawRows(pdf, document.Columns, document.Rows, widths, fontName)

	return pdf.OutputFileAndClose(outputPath)
}

func determineOrientation(columns []model.Column, requested string) string {
	switch strings.ToLower(strings.TrimSpace(requested)) {
	case "portrait":
		return "P"
	case "landscape":
		return "L"
	case "auto", "":
		if len(columns) > 5 {
			return "L"
		}
		return "P"
	default:
		return "P"
	}
}

func resolveFontPath() string {
	candidates := []string{}
	switch runtime.GOOS {
	case "windows":
		candidates = []string{
			`C:\Windows\Fonts\arial.ttf`,
			`C:\Windows\Fonts\calibri.ttf`,
			`C:\Windows\Fonts\tahoma.ttf`,
			`C:\Windows\Fonts\verdana.ttf`,
		}
	case "darwin":
		candidates = []string{
			"/System/Library/Fonts/Supplemental/Arial.ttf",
			"/System/Library/Fonts/Helvetica.ttc",
			"/Library/Fonts/Arial.ttf",
		}
	default:
		candidates = []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/liberation2/LiberationSans-Regular.ttf",
			"/usr/share/fonts/truetype/dejavu/DejaVuSansCondensed.ttf",
		}
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func columnWidths(pdf *gofpdf.Fpdf, count int) []float64 {
	if count == 0 {
		return nil
	}

	pageWidth, _ := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	usableWidth := pageWidth - left - right
	widths := make([]float64, count)
	base := usableWidth / float64(count)
	for i := range widths {
		widths[i] = base
	}
	return widths
}

func drawHeader(pdf *gofpdf.Fpdf, columns []model.Column, widths []float64, fontName string) {
	if len(columns) == 0 {
		return
	}

	pdf.SetFillColor(224, 233, 245)
	pdf.SetTextColor(25, 25, 25)
	pdf.SetFont(fontName, "", 9)
	x := pdf.GetX()
	y := pdf.GetY()

	for i, col := range columns {
		pdf.SetXY(x, y)
		pdf.CellFormat(widths[i], 7, strings.TrimSpace(col.Label), "1", 0, "C", true, 0, "")
		x += widths[i]
	}

	pdf.Ln(7)
}

func drawRows(pdf *gofpdf.Fpdf, columns []model.Column, rows []model.Row, widths []float64, fontName string) {
	if len(rows) == 0 {
		return
	}

	_, pageHeight := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()

	pdf.SetFont(fontName, "", 8)
	pdf.SetTextColor(35, 35, 35)

	for _, row := range rows {
		if pdf.GetY()+10 > pageHeight-bottom-12 {
			pdf.AddPage()
			drawHeader(pdf, columns, widths, fontName)
		}

		x := pdf.GetX()
		y := pdf.GetY()
		for i, col := range columns {
			value := row.Values[col.Key]
			formatted := formatter.FormatValue(value, col.Type)
			pdf.SetXY(x, y)
			pdf.CellFormat(widths[i], 6, formatted, "1", 0, "L", false, 0, "")
			x += widths[i]
		}
		pdf.Ln(6)
	}
}
