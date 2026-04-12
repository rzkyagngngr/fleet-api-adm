package main

import (
	"fmt"
	"path/filepath"

	"omniport-api/cmd/docs/trd/generator"
)

func main() {
	swaggerPath := "docs/swagger.json"
	templatePath := "docs/trd_template.docx"
	outputDir := "docs/trd"

	// Mapping of Module Name to Swagger Tags
	modules := []struct {
		Name string
		Tags []string
	}{
		{"Access", []string{"master-role-access"}},
		{"Auth", []string{"auth"}},
		{"Dermaga", []string{"dermaga"}},
		{"Menu", []string{"master-menus"}},
		{"Reference", []string{"master-references"}},
		{"Role", []string{"master-roles"}},
		{"User", []string{"users"}},
	}

	for _, mod := range modules {
		fmt.Printf("Generating TRD for %s...\n", mod.Name)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("TRD - %s.docx", mod.Name))
		
		err := generator.GenerateTRD(swaggerPath, templatePath, outputPath, mod.Name, mod.Tags)
		if err != nil {
			fmt.Printf("  Error generating %s: %v\n", mod.Name, err)
		} else {
			fmt.Printf("  Successfully generated: %s\n", outputPath)
		}
	}

	fmt.Println("\nTRD generation process completed.")
}
