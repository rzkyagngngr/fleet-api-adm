package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"omniport-api/cmd/docs/trd/generator"
)

func main() {
	swaggerPath := "docs/swagger.json"
	templatePath := "docs/trd_template.docx"
	outputDir := "docs/trd"

	modulesRoot := "internal/modules"
	categories, err := os.ReadDir(modulesRoot)
	if err != nil {
		fmt.Printf("Error reading modules root: %v\n", err)
		return
	}

	for _, cat := range categories {
		if !cat.IsDir() {
			continue
		}
		
		modulesPath := filepath.Join(modulesRoot, cat.Name())
		entries, err := os.ReadDir(modulesPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			
			name := entry.Name()
		displayName := strings.Title(name)
		
		// Map directory name to swagger tags
		tags := []string{name, "master-"+name, "master-"+name+"s"}
		switch name {
		case "access":
			tags = []string{"master-role-access"}
		case "menu":
			tags = []string{"master-menus"}
		case "reference":
			tags = []string{"master-references"}
		case "role":
			tags = []string{"master-roles"}
		case "user":
			tags = []string{"users", "master-users"}
		case "auth":
			tags = []string{"auth"}
		case "company":
			tags = []string{"master-company"}
		case "branch":
			tags = []string{"master-branches", "master-branch"}
		case "terminal":
			tags = []string{"master-terminals", "master-terminal"}
		case "pelabuhan":
			tags = []string{"master-pelabuhan", "master-port", "master-ports"}
		case "customer":
			tags = []string{"master-customer", "master-customers"}
		case "vessel":
			tags = []string{"master-vessel", "master-vessels"}
		case "cargo":
			tags = []string{"master-barang", "barang", "cargo"}
		case "warehouse":
			tags = []string{"master-warehouse", "master-warehouses"}
		case "equipment":
			tags = []string{"master-equipment", "master-equipments"}
		}

		fmt.Printf("Generating TRD for %s (tags: %v)...\n", displayName, tags)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("TRD - %s.docx", displayName))
		
		err := generator.GenerateTRD(swaggerPath, templatePath, outputPath, displayName, tags)
		if err != nil {
			fmt.Printf("  Error generating %s: %v\n", displayName, err)
		} else {
			fmt.Printf("  Successfully generated: %s\n", outputPath)
		}
		}
	}

	fmt.Println("\nTRD generation process completed.")
}
