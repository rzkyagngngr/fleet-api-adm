package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"omniport-api/cmd/docs/pia/generator"
)

func main() {
	cmd := exec.Command(
		"swag",
		"init",
		"-g", "main.go",
		"-d", "cmd/monolith,internal/modules/administration/auth,internal/modules/administration/user,internal/modules/administration/menu,internal/modules/administration/role,internal/modules/administration/access,internal/modules/administration/reference,internal/modules/administration/dermaga,internal/helper,internal/router",
		"--parseInternal",
		"-o", "docs",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "swagger generation failed: %v\n", err)
		os.Exit(1)
	}

	b, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "swagger docs generated but failed to read docs/swagger.json: %v\n", err)
		os.Exit(1)
	}

	if strings.Contains(string(b), `"paths": {}`) {
		fmt.Fprintln(os.Stdout, "swagger generated, but paths are empty.")
		fmt.Fprintln(os.Stdout, "add swag annotations (@Summary, @Tags, @Param, @Success, @Failure, @Router) in handlers.")
	}

	if err := generator.GenerateFromSwaggerAndTemplate("docs/swagger.json", "docs/pia_template.docx", "docs/pia/pia_generated.docx"); err != nil {
		if strings.Contains(err.Error(), "target output is locked; generated fallback file:") {
			fmt.Fprintf(os.Stdout, "swagger generated; %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "swagger generated but PIA generation failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("swagger docs generated successfully")
	fmt.Println("pia docx generated successfully: docs/pia/pia_generated.docx")
}
