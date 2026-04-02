package generator

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/beevik/etree"
)

type swaggerSpec struct {
	BasePath    string                          `json:"basePath"`
	Paths       map[string]map[string]operation `json:"paths"`
	Definitions map[string]schemaDef            `json:"definitions"`
}

type operation struct {
	Summary    string              `json:"summary"`
	Tags       []string            `json:"tags"`
	Consumes   []string            `json:"consumes"`
	Parameters []parameter         `json:"parameters"`
	Responses  map[string]response `json:"responses"`
}

type parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Schema   struct {
		Ref string `json:"$ref"`
	} `json:"schema"`
	Description string `json:"description"`
}

type response struct {
	Description string `json:"description"`
}

type schemaDef struct {
	Type       string                 `json:"type"`
	Properties map[string]propertyDef `json:"properties"`
}

type propertyDef struct {
	Type string `json:"type"`
	Ref  string `json:"$ref"`
}

type endpointEntry struct {
	Method string
	Path   string
	Op     operation
}

func GenerateFromSwaggerAndTemplate(swaggerPath, templatePath, outputPath string) error {
	spec, err := readSwagger(swaggerPath)
	if err != nil {
		return err
	}
	if len(spec.Paths) == 0 {
		return fmt.Errorf("swagger paths are empty in %s", swaggerPath)
	}

	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	docxBytes, err := buildDocxWithTemplateStyle(templateBytes, spec, flattenEndpoints(spec))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	if err := os.WriteFile(outputPath, docxBytes, 0o644); err != nil {
		alt := withTimestampSuffix(outputPath)
		if err2 := os.WriteFile(alt, docxBytes, 0o644); err2 != nil {
			return fmt.Errorf("failed to write output docx: %w (and fallback failed: %v)", err, err2)
		}
		return fmt.Errorf("target output is locked; generated fallback file: %s", alt)
	}
	return nil
}

func readSwagger(path string) (*swaggerSpec, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read swagger json: %w", err)
	}
	var spec swaggerSpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse swagger json: %w", err)
	}
	return &spec, nil
}

func flattenEndpoints(spec *swaggerSpec) []endpointEntry {
	out := make([]endpointEntry, 0, 64)
	for p, methods := range spec.Paths {
		for m, op := range methods {
			out = append(out, endpointEntry{
				Method: strings.ToUpper(m),
				Path:   p,
				Op:     op,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return methodOrder(out[i].Method) < methodOrder(out[j].Method)
	})
	return out
}

func methodOrder(m string) int {
	switch m {
	case "GET":
		return 1
	case "POST":
		return 2
	case "PUT":
		return 3
	case "PATCH":
		return 4
	case "DELETE":
		return 5
	default:
		return 99
	}
}

func buildDocxWithTemplateStyle(templateBytes []byte, spec *swaggerSpec, endpoints []endpointEntry) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(templateBytes), int64(len(templateBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open template docx: %w", err)
	}

	var out bytes.Buffer
	w := zip.NewWriter(&out)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("failed to open zip entry %s: %w", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("failed to read zip entry %s: %w", f.Name, err)
		}

		if f.Name == "word/document.xml" {
			data, err = buildStyledDocumentXML(data, spec, endpoints)
			if err != nil {
				_ = w.Close()
				return nil, err
			}
		}

		h := f.FileHeader
		dst, err := w.CreateHeader(&h)
		if err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("failed to create output zip entry %s: %w", f.Name, err)
		}
		if _, err := dst.Write(data); err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("failed to write output zip entry %s: %w", f.Name, err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize output docx: %w", err)
	}
	return out.Bytes(), nil
}

func buildStyledDocumentXML(docXML []byte, spec *swaggerSpec, endpoints []endpointEntry) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(docXML); err != nil {
		return nil, fmt.Errorf("failed to parse document.xml: %w", err)
	}

	body := doc.FindElement("//w:body")
	if body == nil {
		body = doc.FindElement("//body")
	}
	if body == nil {
		return nil, fmt.Errorf("invalid document.xml: w:body not found")
	}

	children := body.ChildElements()
	start := -1
	end := -1
	methodIdx := -1
	responseIdx := -1
	for i, c := range children {
		if isTag(c, "tbl") && methodIdx == -1 && hasCellText(c, "METHOD") && hasCellText(c, "URL") {
			methodIdx = i
		}
		if isTag(c, "tbl") && hasCellText(c, "STATUS") && hasCellText(c, "RESPONSE") {
			responseIdx = i
		}
	}
	if methodIdx >= 0 {
		start = methodIdx
		for j := methodIdx - 1; j >= 0; j-- {
			if isTag(children[j], "p") && containsText(children[j], "Request") {
				start = j
				break
			}
			if isTag(children[j], "tbl") {
				break
			}
		}
	}
	if responseIdx >= 0 {
		end = responseIdx
	}
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("failed to locate template request/response block in document")
	}

	templateBlock := make([]*etree.Element, 0, end-start+1)
	for i := start; i <= end; i++ {
		templateBlock = append(templateBlock, children[i].Copy())
	}

	for i := end; i >= start; i-- {
		body.RemoveChild(children[i])
	}

	insertAt := start
	for _, ep := range endpoints {
		clones := cloneBlock(templateBlock)
		applyEndpointToBlock(clones, spec, ep)
		for _, n := range clones {
			anchor := body.ChildElements()
			if insertAt >= len(anchor) {
				body.AddChild(n)
			} else {
				body.InsertChildAt(insertAt, n)
			}
			insertAt++
		}
		body.InsertChildAt(insertAt, cloneParagraph(templateBlock[0]))
		insertAt++
	}

	doc.Indent(0)
	out, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize document.xml: %w", err)
	}
	return out, nil
}

func cloneBlock(block []*etree.Element) []*etree.Element {
	out := make([]*etree.Element, 0, len(block))
	for _, b := range block {
		out = append(out, b.Copy())
	}
	return out
}

func cloneParagraph(from *etree.Element) *etree.Element {
	return from.Copy()
}

func applyEndpointToBlock(block []*etree.Element, spec *swaggerSpec, ep endpointEntry) {
	tables := make([]*etree.Element, 0, 5)
	for _, b := range block {
		if isTag(b, "tbl") {
			tables = append(tables, b)
		}
	}
	if len(tables) < 5 {
		return
	}

	url := joinPath(spec.BasePath, ep.Path)
	setCellText(tables[0], 1, 0, ep.Method)
	setCellText(tables[0], 1, 1, url)

	setCellText(tables[1], 0, 1, buildRequestExample(ep, spec))
	setCellText(tables[1], 1, 1, inferPrimaryKey(ep))

	setCellText(tables[2], 1, 0, inferPayloadType(ep))
	setCellText(tables[2], 1, 1, buildParams(ep, "name"))
	setCellText(tables[2], 1, 2, buildParams(ep, "desc"))
	setCellText(tables[2], 1, 3, buildParams(ep, "type"))

	setCellText(tables[3], 0, 1, buildHeaders(ep))

	setCellText(tables[4], 1, 0, buildStatusCodes(ep))
	setCellText(tables[4], 1, 1, buildResponseSummary(ep))
}

func setCellText(tbl *etree.Element, rowIdx, colIdx int, value string) {
	rows := findDesc(tbl, "tr")
	if rowIdx < 0 || rowIdx >= len(rows) {
		return
	}
	cols := make([]*etree.Element, 0, 4)
	for _, c := range rows[rowIdx].ChildElements() {
		if isTag(c, "tc") {
			cols = append(cols, c)
		}
	}
	if colIdx < 0 || colIdx >= len(cols) {
		return
	}
	tc := cols[colIdx]
	baseP := firstDesc(tc, "p")
	if baseP == nil {
		return
	}

	lines := splitLines(value)
	if len(lines) == 0 {
		lines = []string{""}
	}

	// Keep only tcPr and rebuild content paragraphs from scratch so old
	// blank placeholders in template do not create huge vertical empty gaps.
	children := tc.ChildElements()
	toRemove := make([]*etree.Element, 0, len(children))
	var tcPr *etree.Element
	for _, ch := range children {
		if isTag(ch, "tcPr") {
			tcPr = ch
			continue
		}
		toRemove = append(toRemove, ch)
	}
	for _, ch := range toRemove {
		tc.RemoveChild(ch)
	}
	if tcPr == nil {
		tcPr = etree.NewElement("w:tcPr")
		tc.InsertChildAt(0, tcPr)
	}

	for _, line := range lines {
		p := baseP.Copy()
		texts := findDesc(p, "t")
		if len(texts) == 0 {
			// fallback paragraph with simple run/text
			p = etree.NewElement("w:p")
			r := etree.NewElement("w:r")
			t := etree.NewElement("w:t")
			t.SetText(line)
			r.AddChild(t)
			p.AddChild(r)
			tc.AddChild(p)
			continue
		}
		texts[0].SetText(line)
		for i := 1; i < len(texts); i++ {
			texts[i].SetText("")
		}
		tc.AddChild(p)
	}
}

func splitLines(s string) []string {
	raw := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		out = append(out, v)
	}
	return out
}

func hasExactText(el *etree.Element, txt string) bool {
	for _, t := range findDesc(el, "t") {
		if strings.TrimSpace(t.Text()) == txt {
			return true
		}
	}
	return false
}

func containsText(el *etree.Element, txt string) bool {
	target := strings.ToUpper(strings.TrimSpace(txt))
	for _, t := range findDesc(el, "t") {
		if strings.Contains(strings.ToUpper(strings.TrimSpace(t.Text())), target) {
			return true
		}
	}
	return false
}

func hasCellText(tbl *etree.Element, txt string) bool {
	target := strings.ToUpper(strings.TrimSpace(txt))
	for _, t := range findDesc(tbl, "t") {
		if strings.ToUpper(strings.TrimSpace(t.Text())) == target {
			return true
		}
	}
	return false
}

func joinPath(base, p string) string {
	b := strings.TrimSpace(base)
	if b == "" {
		return p
	}
	if strings.HasSuffix(b, "/") {
		b = strings.TrimSuffix(b, "/")
	}
	return b + p
}

func buildRequestExample(ep endpointEntry, spec *swaggerSpec) string {
	bodyParams := make([]parameter, 0, 1)
	for _, p := range ep.Op.Parameters {
		if p.In == "body" {
			bodyParams = append(bodyParams, p)
		}
	}
	if len(bodyParams) == 0 {
		return "{}"
	}
	ref := bodyParams[0].Schema.Ref
	if ref == "" {
		return "{}"
	}

	schemaName := refName(ref)
	schema, ok := spec.Definitions[schemaName]
	if !ok || len(schema.Properties) == 0 {
		return "{\n  \"payload\": \"" + schemaName + "\"\n}"
	}

	keys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys)+2)
	lines = append(lines, "{")
	for i, k := range keys {
		p := schema.Properties[k]
		v := sampleValue(p.Type)
		suffix := ","
		if i == len(keys)-1 {
			suffix = ""
		}
		lines = append(lines, fmt.Sprintf("  \"%s\": %s%s", k, v, suffix))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func sampleValue(typ string) string {
	switch strings.ToLower(strings.TrimSpace(typ)) {
	case "string":
		return "\"string\""
	case "integer", "number":
		return "0"
	case "boolean":
		return "true"
	case "array":
		return "[]"
	case "object":
		return "{}"
	default:
		return "\"value\""
	}
}

func inferPrimaryKey(ep endpointEntry) string {
	for _, p := range ep.Op.Parameters {
		if p.In == "path" && strings.TrimSpace(p.Name) != "" {
			return p.Name
		}
	}
	for _, p := range ep.Op.Parameters {
		if p.Required {
			return p.Name
		}
	}
	return "-"
}

func inferPayloadType(ep endpointEntry) string {
	if len(ep.Op.Consumes) > 0 {
		if strings.Contains(strings.ToLower(ep.Op.Consumes[0]), "json") {
			return "JSON"
		}
		return ep.Op.Consumes[0]
	}
	return "JSON"
}

func buildParams(ep endpointEntry, mode string) string {
	if len(ep.Op.Parameters) == 0 {
		return "-"
	}
	lines := make([]string, 0, len(ep.Op.Parameters))
	for _, p := range ep.Op.Parameters {
		switch mode {
		case "name":
			lines = append(lines, p.Name)
		case "desc":
			if strings.TrimSpace(p.Description) == "" {
				lines = append(lines, "Parameter "+p.Name)
			} else {
				lines = append(lines, p.Description)
			}
		case "type":
			t := p.Type
			if t == "" && p.Schema.Ref != "" {
				t = refName(p.Schema.Ref)
			}
			if t == "" {
				t = "-"
			}
			lines = append(lines, t)
		}
	}
	return strings.Join(lines, "\n")
}

func buildHeaders(ep endpointEntry) string {
	if len(ep.Op.Consumes) == 0 {
		return "Content-type: application/json"
	}
	return "Content-type: " + ep.Op.Consumes[0]
}

func buildStatusCodes(ep endpointEntry) string {
	if len(ep.Op.Responses) == 0 {
		return "-"
	}
	codes := make([]string, 0, len(ep.Op.Responses))
	for c := range ep.Op.Responses {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	return strings.Join(codes, "\n")
}

func buildResponseSummary(ep endpointEntry) string {
	if len(ep.Op.Responses) == 0 {
		return "{}"
	}
	codes := make([]string, 0, len(ep.Op.Responses))
	for c := range ep.Op.Responses {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	lines := make([]string, 0, len(codes)+2)
	lines = append(lines, "{")
	for _, c := range codes {
		desc := strings.TrimSpace(ep.Op.Responses[c].Description)
		if desc == "" {
			desc = "Response"
		}
		lines = append(lines, fmt.Sprintf("  \"%s\": \"%s\",", c, desc))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func refName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

func withTimestampSuffix(path string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	return fmt.Sprintf("%s_%s%s", base, time.Now().Format("20060102_150405"), ext)
}

func isTag(el *etree.Element, local string) bool {
	if el == nil {
		return false
	}
	return localName(el.Tag) == local
}

func localName(tag string) string {
	if idx := strings.Index(tag, ":"); idx >= 0 {
		return tag[idx+1:]
	}
	return tag
}

func findDesc(root *etree.Element, local string) []*etree.Element {
	out := make([]*etree.Element, 0, 8)
	var walk func(n *etree.Element)
	walk = func(n *etree.Element) {
		for _, c := range n.ChildElements() {
			if isTag(c, local) {
				out = append(out, c)
			}
			walk(c)
		}
	}
	walk(root)
	return out
}

func firstDesc(root *etree.Element, local string) *etree.Element {
	var found *etree.Element
	var walk func(n *etree.Element)
	walk = func(n *etree.Element) {
		if found != nil {
			return
		}
		for _, c := range n.ChildElements() {
			if isTag(c, local) {
				found = c
				return
			}
			walk(c)
			if found != nil {
				return
			}
		}
	}
	walk(root)
	return found
}
