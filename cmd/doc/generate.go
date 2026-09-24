package doc

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"blexios/internal/doc/service"
)

func RunGenerate(args []string) int {
	if len(args) == 0 {
		printRootHelp()
		return 0
	}

	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		printRootHelp()
		return 0
	}

	if args[0] != "doc:generate" {
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printRootHelp()
		return 1
	}

	fs := flag.NewFlagSet("doc:generate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	opts := service.Options{}
	fs.BoolVar(&opts.PDF, "pdf", false, "Generate PDF")
	fs.StringVar(&opts.Path, "path", "", "Input JSON file")
	fs.StringVar(&opts.Output, "output", "", "Output file")
	fs.StringVar(&opts.Title, "title", "", "Document title")
	fs.StringVar(&opts.Template, "template", "default", "Document template")
	fs.StringVar(&opts.Orientation, "orientation", "auto", "portrait, landscape or auto")
	fs.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose output")

	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			printGenerateHelp()
			return 0
		}
		return 1
	}

	if isHelpFlagPresent(args[1:]) {
		printGenerateHelp()
		return 0
	}

	if !opts.PDF {
		fmt.Println("Error: a format must be selected.")
		printGenerateHelp()
		return 1
	}

	if opts.Path == "" {
		fmt.Println("Error: input file not found:")
		fmt.Printf("%s\n", opts.Path)
		return 1
	}

	generator := service.NewGenerator()
	doc, outputPath, err := generator.Generate(opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	fmt.Println("Blexios DOC Generator")
	fmt.Println()
	fmt.Println("Input:")
	fmt.Println(opts.Path)
	fmt.Printf("Records: %d\n", doc.Metadata.RecordCount)
	fmt.Printf("Columns: %d\n", doc.Metadata.ColumnCount)
	fmt.Println()
	fmt.Println("Generating PDF...")
	fmt.Println("████████████████████████ 100%")
	fmt.Println()
	fmt.Println("PDF generated successfully.")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println(outputPath)

	return 0
}

func isHelpFlagPresent(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func printRootHelp() {
	fmt.Println("Blexios Document Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  blexios doc:generate [flags]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  blexios doc:generate --pdf --path=data.json")
	fmt.Println("  blexios doc:generate --pdf --path=data.json --output=rapport.pdf")
	fmt.Println()
	printGenerateHelp()
}

func printGenerateHelp() {
	fmt.Println("Generate documents from structured data.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  blexios doc:generate [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --pdf                  Generate PDF")
	fmt.Println("  --path string          Input JSON file")
	fmt.Println("  --output string        Output file")
	fmt.Println("  --title string         Document title")
	fmt.Println("  --template string      Document template")
	fmt.Println("  --orientation string   portrait, landscape or auto")
	fmt.Println("  --verbose              Enable verbose output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  blexios doc:generate --pdf --path=data.json")
	fmt.Println("  blexios doc:generate --pdf --path=data.json --output=rapport.pdf")
	fmt.Println("  blexios doc:generate --pdf --path=data.json --title=\"Rapport de présence\"")
}

func normalizeTitle(title string) string {
	if title == "" {
		return "Blexios Document"
	}
	return strings.TrimSpace(title)
}
