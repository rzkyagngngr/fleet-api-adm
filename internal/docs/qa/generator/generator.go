package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type TestCase struct {
	Module           string
	Name             string
	Description      string
	Action           string
	ExpectedResult   string
	Endpoint         string
	Method           string
	Payload          string
	ExpectedResponse string
}

func Generate(sprint string, modules []string, templatePath string, outputPath string) error {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}
	defer f.Close()

	sheetName := "Sheet1"
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		sheetName = f.GetSheetList()[0]
	}

	row := 6
	devPrefix := "http://dev-tuks-api.ilcs.co.id"
	sysdate := time.Now().Format("2006-01-02")

	// Styles
	positiveStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#B7E1CD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	moduleHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	totalTestCases := 0

	for _, mod := range modules {
		// Module Header
		f.MergeCell(sheetName, fmt.Sprintf("A%d", row-1), fmt.Sprintf("T%d", row-1))
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row-1), "MODULE: "+strings.ToUpper(mod))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row-1), fmt.Sprintf("T%d", row-1), moduleHeaderStyle)

		testCases := generateTestCasesForModule(mod)
		for i, tc := range testCases {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tc.Module)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tc.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "API")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), sprint)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "positif")
			f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), positiveStyle)

			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), normalizeDescription(tc.Name))
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), fmt.Sprintf("%s %s", tc.Method, tc.Endpoint))
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), "success")
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), 0) // Default to 0
			f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), 0) // Default to 0
			
			// N-T Columns
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), "Sofyan")
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), sysdate)
			f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), devPrefix+tc.Endpoint)
			f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), tc.Payload)
			f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), tc.ExpectedResponse)
			f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), generateCurl(tc.Method, devPrefix+tc.Endpoint, tc.Payload))

			// Apply styles
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), dataStyle)
			f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("T%d", row), dataStyle)

			totalTestCases++
			row++
		}
		row += 2
	}

	// Summary Statistics with Formulas
	lastRow := row - 3 // Adjust for the last gap and header
	if totalTestCases > 0 {
		// J1: Percentage (Passed / Total)
		// We use COUNT to get total rows and SUM for passed
		f.SetCellFormula(sheetName, "J1", fmt.Sprintf("IFERROR(SUM(K6:K%d)/COUNT(A6:A%d)*100, 0)", lastRow, lastRow))
		
		// K1: [SUM]/[TOTAL] Test Passed
		// Concatenate sum and total
		f.SetCellFormula(sheetName, "K1", fmt.Sprintf("CONCATENATE(SUM(K6:K%d), \"/\", COUNT(A6:A%d), \" Test Passed\")", lastRow, lastRow))
		
		// K2: [SUM]/[TOTAL] Test Failed
		f.SetCellFormula(sheetName, "K2", fmt.Sprintf("CONCATENATE(SUM(L6:L%d), \"/\", COUNT(A6:A%d), \" Test Failed\")", lastRow, lastRow))
	}
	
	// Add styling for summary
	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(sheetName, "J1", "K2", summaryStyle)

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	if err := f.SaveAs(outputPath); err != nil {
		ext := filepath.Ext(outputPath)
		base := strings.TrimSuffix(outputPath, ext)
		altPath := fmt.Sprintf("%s_%s%s", base, time.Now().Format("20060102_150405"), ext)
		if err2 := f.SaveAs(altPath); err2 != nil {
			return fmt.Errorf("failed to save excel: %w (and fallback failed: %v)", err, err2)
		}
		fmt.Printf("Warning: Target file was locked. Generated fallback: %s\n", altPath)
	}

	return nil
}

func normalizeDescription(name string) string {
	words := strings.Fields(name)
	if len(words) == 0 {
		return ""
	}
	verb := strings.ToLower(words[0])
	rest := strings.Join(words[1:], " ")
	
	switch verb {
	case "get":
		verb = "getting"
	case "create":
		verb = "creating"
	case "update":
		verb = "updating"
	case "delete":
		verb = "deleting"
	case "search":
		verb = "searching"
	case "save":
		verb = "saving"
	}
	
	return fmt.Sprintf("%s %s", verb, rest)
}

func generateCurl(method, url, payload string) string {
	if payload == "" || payload == "-" {
		return fmt.Sprintf("curl -X %s \"%s\"", method, url)
	}
	// Clean payload for curl
	p := strings.ReplaceAll(payload, "\n", "")
	p = strings.ReplaceAll(p, "\"", "\\\"")
	return fmt.Sprintf("curl -X %s \"%s\" -H \"Content-Type: application/json\" -d \"%s\"", method, url, p)
}

func generateTestCasesForModule(module string) []TestCase {
	var tcs []TestCase

	switch module {
	case "auth":
		tcs = append(tcs, 
			TestCase{Module: module, Name: "User Registration", Method: "POST", Endpoint: "/api/v1/auth/register", Payload: `{"email":"test@example.com","password":"password123"}`, ExpectedResponse: `{"status":201,"message":"Created"}`},
			TestCase{Module: module, Name: "User Login", Method: "POST", Endpoint: "/api/v1/auth/login", Payload: `{"username":"admin","password":"password"}`, ExpectedResponse: `{"status":200,"message":"OK"}`},
		)
	case "user":
		tcs = append(tcs, 
			TestCase{Module: module, Name: "Get Profile", Method: "GET", Endpoint: "/api/v1/users/profile", Payload: "-", ExpectedResponse: `{"status":200}`},
			TestCase{Module: module, Name: "Search Users", Method: "POST", Endpoint: "/api/v1/master/users/search", Payload: `{"filter":{}}`, ExpectedResponse: `{"status":200}`},
		)
	case "menu":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/menus")...)
	case "role":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/roles")...)
	case "access":
		tcs = append(tcs, 
			TestCase{Module: module, Name: "Get Role Access", Method: "GET", Endpoint: "/api/v1/master/roles/1/access", Payload: "-", ExpectedResponse: `{"status":200}`},
			TestCase{Module: module, Name: "Update Role Access", Method: "POST", Endpoint: "/api/v1/master/roles/1/access", Payload: `{"menu_ids":[1,2]}`, ExpectedResponse: `{"status":200}`},
		)
	case "company":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/company")...)
	case "branch":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/branches")...)
	case "terminal":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/terminals")...)
	case "pelabuhan":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/pelabuhan")...)
	case "dermaga":
		tcs = append(tcs, generateCRUD(module, "/api/v1/dermaga")...)
	case "customer":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/customer")...)
	case "vessel":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/vessel")...)
	case "cargo":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/barang")...)
	case "warehouse":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/warehouse")...)
	case "equipment":
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/equipment")...)
	default:
		tcs = append(tcs, generateCRUD(module, "/api/v1/master/"+module)...)
	}

	return tcs
}

func generateCRUD(module, path string) []TestCase {
	baseDesc := strings.Title(module)
	return []TestCase{
		{Module: module, Name: "Search " + baseDesc, Method: "POST", Endpoint: path + "/search", Payload: `{"filter":{}}`, ExpectedResponse: `{"status":200,"data":[]}`},
		{Module: module, Name: "Create " + baseDesc, Method: "POST", Endpoint: path, Payload: `{"name":"test"}`, ExpectedResponse: `{"status":201}`},
		{Module: module, Name: "Update " + baseDesc, Method: "PUT", Endpoint: path + "/1", Payload: `{"name":"updated"}`, ExpectedResponse: `{"status":200}`},
		{Module: module, Name: "Delete " + baseDesc, Method: "DELETE", Endpoint: path + "/1", Payload: "-", ExpectedResponse: `{"status":200}`},
	}
}
