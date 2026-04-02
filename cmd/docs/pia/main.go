package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"omniport-api/cmd/docs/pia/generator"
)

func main() {
	swaggerPath := flag.String("swagger", "docs/swagger.json", "Path to generated swagger JSON")
	templatePath := flag.String("template", "docs/pia_template.docx", "Path to PIA DOCX template")
	outputPath := flag.String("out", "docs/pia/pia_generated.docx", "Output DOCX path")
	flag.Parse()

	if err := generator.GenerateFromSwaggerAndTemplate(*swaggerPath, *templatePath, *outputPath); err != nil {
		if strings.Contains(err.Error(), "target output is locked; generated fallback file:") {
			fmt.Fprintf(os.Stdout, "PIA generation warning: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "PIA generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "PIA docx generated successfully: %s\n", *outputPath)
}
