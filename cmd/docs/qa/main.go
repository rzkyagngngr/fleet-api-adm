package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"omniport-api/internal/docs/qa/generator"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== QA Test Script Generator (Fleet API) ===")
	fmt.Print("Sprint Number: ")
	sprint, _ := reader.ReadString('\n')
	sprint = strings.TrimSpace(sprint)

	modulesRoot := "internal/modules"
	categories, _ := os.ReadDir(modulesRoot)
	var modules []string
	for _, cat := range categories {
		if cat.IsDir() {
			subModules, _ := listModules(filepath.Join(modulesRoot, cat.Name()))
			modules = append(modules, subModules...)
		}
	}

	fmt.Println("\nAvailable Modules:")
	for i, m := range modules {
		fmt.Printf("[%d] %s\n", i+1, m)
	}

	fmt.Print("\nSelect modules (e.g., 1,2,5 or 'all'): ")
	selection, _ := reader.ReadString('\n')
	selection = strings.TrimSpace(selection)

	var selectedModules []string
	if strings.ToLower(selection) == "all" {
		selectedModules = modules
	} else {
		parts := strings.Split(selection, ",")
		for _, p := range parts {
			var idx int
			fmt.Sscanf(strings.TrimSpace(p), "%d", &idx)
			if idx > 0 && idx <= len(modules) {
				selectedModules = append(selectedModules, modules[idx-1])
			}
		}
	}

	if len(selectedModules) == 0 {
		fmt.Println("No modules selected. Exiting.")
		return
	}

	fmt.Printf("\nGenerating Excel for Sprint %s with modules: %v...\n", sprint, selectedModules)
	
	templatePath := "docs/qa_test_template.xlsx"
	outputPath := fmt.Sprintf("docs/qa/Sprint%s_Administration_Tests.xlsx", sprint)

	err := generator.Generate(sprint, selectedModules, templatePath, outputPath)
	if err != nil {
		fmt.Printf("Error generating Excel: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated: %s\n", outputPath)
}

func listModules(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, entry := range entries {
		if entry.IsDir() {
			modules = append(modules, entry.Name())
		}
	}
	return modules, nil
}
