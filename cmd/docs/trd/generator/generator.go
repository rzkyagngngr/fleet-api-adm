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
	Responses  map[string]responses `json:"responses"`
}

type responses struct {
	Description string `json:"description"`
}

type parameter struct {
	Name        string     `json:"name"`
	In          string     `json:"in"`
	Required    bool       `json:"required"`
	Type        string     `json:"type"`
	Schema      *schemaDef `json:"schema"`
	Description string     `json:"description"`
}

type schemaDef struct {
	Type       string                 `json:"type"`
	Ref        string                 `json:"$ref"`
	Items      *propertyDef           `json:"items"`
	Properties map[string]propertyDef `json:"properties"`
}

type propertyDef struct {
	Type   string                 `json:"type"`
	Ref    string                 `json:"$ref"`
	Items  *propertyDef           `json:"items"`
	Format string                 `json:"format"`
	Props  map[string]propertyDef `json:"properties"`
}

type endpointEntry struct {
	Method string
	Path   string
	Op     operation
}

func GenerateTRD(swaggerPath, templatePath, outputPath, moduleName string, targetTags []string) error {
	spec, err := readSwagger(swaggerPath)
	if err != nil { return err }
	endpoints := filterEndpoints(spec, targetTags)
	if len(endpoints) == 0 { return fmt.Errorf("no endpoints found") }
	
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil { return err }
	docxBytes, err2 := buildTRDDocx(templateBytes, spec, endpoints, moduleName)
	if err2 != nil { return err2 }
	
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil { return err }
	return os.WriteFile(outputPath, docxBytes, 0o644)
}

func readSwagger(path string) (*swaggerSpec, error) {
	b, err := os.ReadFile(path)
	if err != nil { return nil, err }
	var spec swaggerSpec
	json.Unmarshal(b, &spec)
	return &spec, nil
}

func filterEndpoints(spec *swaggerSpec, tags []string) []endpointEntry {
	out := make([]endpointEntry, 0)
	tm := make(map[string]bool)
	for _, t := range tags { tm[t] = true }
	for p, ms := range spec.Paths {
		for m, op := range ms {
			match := false
			for _, t := range op.Tags { if tm[t] { match = true; break } }
			if match { out = append(out, endpointEntry{Method: strings.ToUpper(m), Path: p, Op: op}) }
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path { return out[i].Path < out[j].Path }
		return out[i].Method < out[j].Method
	})
	return out
}

func buildTRDDocx(templateBytes []byte, spec *swaggerSpec, endpoints []endpointEntry, moduleName string) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(templateBytes), int64(len(templateBytes)))
	if err != nil { return nil, err }
	var out bytes.Buffer
	w := zip.NewWriter(&out)
	for _, f := range r.File {
		rc, _ := f.Open()
		data, _ := io.ReadAll(rc)
		rc.Close()
		if f.Name == "word/document.xml" {
			data, err = processTRDXML(data, spec, endpoints, moduleName)
			if err != nil { _ = w.Close(); return nil, err }
		}
		dst, _ := w.Create(f.Name)
		dst.Write(data)
	}
	_ = w.Close()
	return out.Bytes(), nil
}

func processTRDXML(docXML []byte, spec *swaggerSpec, endpoints []endpointEntry, moduleName string) ([]byte, error) {
	doc := etree.NewDocument(); doc.ReadFromBytes(docXML)
	body := doc.FindElement("//w:body")
	children := body.ChildElements()

	var tblTemplate *etree.Element
	var pTemplateEndpoint *etree.Element
	templateMarkerIdx := -1
	fiturHeaderIdx := -1

	for i, c := range children {
		txt := getElementText(c)
		txtUpper := strings.ToUpper(txt)
		
		if i == 0 && isTag(c, "p") { ensureParagraphStyle(c, "1", "0") }
		if i == 0 && isTag(c, "tbl") { setCellTextInTable(c, 1, 1, "Omniport | "+moduleName) }

		if strings.Contains(txt, "Deskripsi") || strings.Contains(txt, "Prd") || strings.Contains(txt, "Pia") || strings.Contains(txt, "Fitur") {
			ensureParagraphStyle(c, "2", "1")
		}

		if strings.Contains(txt, "[template_endpoint]") { 
			templateMarkerIdx = i
			pTemplateEndpoint = c.Copy()
		}
		if strings.Contains(txtUpper, "FITUR") && fiturHeaderIdx == -1 { fiturHeaderIdx = i }
		if isTag(c, "tbl") && tblTemplate == nil {
			if strings.Contains(txtUpper, "URL") && strings.Contains(txtUpper, "METHOD") && strings.Contains(txtUpper, "REQUEST BODY") {
				tblTemplate = c.Copy()
			}
		}
	}

	if tblTemplate == nil || pTemplateEndpoint == nil { return nil, fmt.Errorf("missing template elements") }

	toRemove := make([]*etree.Element, 0)
	for i, c := range children {
		if isTag(c, "tbl") && i != 0 {
			txt := strings.ToUpper(getElementText(c))
			if strings.Contains(txt, "URL") && strings.Contains(txt, "METHOD") { toRemove = append(toRemove, c) }
		}
		if i == templateMarkerIdx { toRemove = append(toRemove, c) }
	}
	for _, r := range toRemove { body.RemoveChild(r) }

	if fiturHeaderIdx != -1 {
		insertAt := fiturHeaderIdx + 1
		for i, ep := range endpoints {
			p := etree.NewElement("w:p")
			r := etree.NewElement("w:r"); t := etree.NewElement("w:t"); t.SetText(fmt.Sprintf("%d. %s %s", i+1, ep.Method, ep.Path))
			r.AddChild(t); p.AddChild(r); body.InsertChildAt(insertAt, p); insertAt++
		}
	}

	insertIdx := -1
	for i, c := range body.ChildElements() { if isTag(c, "sectPr") { insertIdx = i; break } }

	bookmarkPrefix := strings.ToLower(moduleName) + "_"
	for i, ep := range endpoints {
		bkName := fmt.Sprintf("%s%d", bookmarkPrefix, i)
		title := fmt.Sprintf("%s %s", ep.Method, ep.Path)
		pTitle := pTemplateEndpoint.Copy()
		setParagraphText(pTitle, title)
		ensureParagraphStyle(pTitle, "3", "2")
		addBookmarkToParagraph(pTitle, 10000+i, bkName)
		
		pTable := tblTemplate.Copy()
		applyEndpointToTRDTable(pTable, spec, ep)
		
		if insertIdx != -1 {
			body.InsertChildAt(insertIdx, pTitle); insertIdx++
			body.InsertChildAt(insertIdx, pTable); insertIdx++
			body.InsertChildAt(insertIdx, etree.NewElement("w:p")); insertIdx++
		} else {
			body.AddChild(pTitle); body.AddChild(pTable); body.AddChild(etree.NewElement("w:p"))
		}
	}
	return doc.WriteToBytes()
}

func ensureParagraphStyle(p *etree.Element, styleId, outlineLvl string) {
	pPr := firstDesc(p, "pPr")
	if pPr == nil { pPr = etree.NewElement("w:pPr"); p.InsertChildAt(0, pPr) }
	pStyle := firstDesc(pPr, "pStyle")
	if pStyle == nil { pStyle = etree.NewElement("w:pStyle"); pPr.InsertChildAt(0, pStyle) }
	pStyle.CreateAttr("w:val", styleId)
	outline := firstDesc(pPr, "outlineLvl")
	if outline == nil { outline = etree.NewElement("w:outlineLvl"); pPr.AddChild(outline) }
	outline.CreateAttr("w:val", outlineLvl)
}

func getElementText(el *etree.Element) string {
	res := ""
	for _, t := range findDesc(el, "t") { res += t.Text() }
	return res
}

func applyEndpointToTRDTable(tbl *etree.Element, spec *swaggerSpec, ep endpointEntry) {
	rows := findDescTopOnly(tbl, "tr")
	toRemove := make([]*etree.Element, 0)
	inHeader := false
	for i := 0; i < len(rows); i++ {
		txt := strings.ToUpper(getElementText(rows[i]))
		if strings.Contains(txt, "REQUEST HEADER") { inHeader = true }
		if strings.Contains(txt, "REQUEST PARAMETER") { inHeader = false; toRemove = append(toRemove, rows[i]) }
		if strings.HasPrefix(txt, "URL") { setCellTextInRow(rows[i], 1, joinPath(spec.BasePath, ep.Path)) }
		if strings.HasPrefix(txt, "METHOD") { setCellTextInRow(rows[i], 1, ep.Method) }
		if inHeader && (strings.Contains(txt, "[HEADER FIELD") || strings.Contains(txt, "AUTHORIZATION")) {
			if !strings.Contains(txt, "[HEADER FIELD 2]") {
				setCellTextInRow(rows[i], 0, "Authorization"); setCellTextInRow(rows[i], 1, "Yes"); setCellTextInRow(rows[i], 2, "Bearer [JWT Token]")
			} else { toRemove = append(toRemove, rows[i]) }
		} else if strings.Contains(txt, "REQUEST BODY") {
			if i+1 < len(rows) {
				if tc := getCell(rows[i+1], 0); tc != nil {
					if inner := firstDesc(tc, "tbl"); inner != nil {
						setCellTextMultiPara(inner, 0, 0, "JSON\nkode\n"+buildRequestExample(ep, spec))
					}
				}
			}
		} else if strings.Contains(txt, "RESPONSE") && !strings.Contains(txt, "SUCCESS") {
			if i+2 < len(rows) {
				rv := rows[i+2]
				if tcS := getCell(rv, 0); tcS != nil {
					if inS := firstDesc(tcS, "tbl"); inS != nil { setCellTextMultiPara(inS, 0, 0, "Plain Text\nkode\n"+buildResponseExample(ep, true)) }
				}
				if tcE := getCell(rv, 2); tcE != nil {
					if inE := firstDesc(tcE, "tbl"); inE != nil { setCellTextMultiPara(inE, 0, 0, "Plain Text\n"+buildResponseExample(ep, false)) }
				}
			}
		} else if strings.Contains(txt, "[PARAMETER_NAME") || strings.Contains(txt, "...") { toRemove = append(toRemove, rows[i]) }
	}
	for _, r := range toRemove { tbl.RemoveChild(r) }
}

func findDescTopOnly(r *etree.Element, l string) []*etree.Element {
	out := []*etree.Element{}
	for _, c := range r.ChildElements() { if isTag(c, l) { out = append(out, c) } }
	return out
}

func addBookmarkToParagraph(p *etree.Element, id int, name string) {
	bs := etree.NewElement("w:bookmarkStart"); bs.CreateAttr("w:id", fmt.Sprintf("%d", id)); bs.CreateAttr("w:name", name)
	be := etree.NewElement("w:bookmarkEnd"); be.CreateAttr("w:id", fmt.Sprintf("%d", id))
	runs := findDesc(p, "r"); if len(runs) > 0 { p.InsertChildAt(indexOfChild(p, runs[0]), bs); p.AddChild(be) }
}

func indexOfChild(p *etree.Element, child *etree.Element) int {
	for i, c := range p.ChildElements() { if c == child { return i } }
	return -1
}

func setCellTextInTable(tbl *etree.Element, r, c int, val string) {
	rows := findDescTopOnly(tbl, "tr")
	if r >= 0 && r < len(rows) { if tc := getCell(rows[r], c); tc != nil { setCellTextSimple(tc, val) } }
}

func setCellTextInRow(row *etree.Element, colIdx int, text string) {
	if tc := getCell(row, colIdx); tc != nil { setCellTextSimple(tc, text) }
}

func setCellTextSimple(tc *etree.Element, value string) {
	tcPr := firstDesc(tc, "tcPr"); baseP := firstDesc(tc, "p")
	if baseP == nil { baseP = etree.NewElement("w:p") }
	for _, ch := range tc.ChildElements() { if !isTag(ch, "tcPr") { tc.RemoveChild(ch) } }
	if tcPr == nil { tcPr = etree.NewElement("w:tcPr"); tc.InsertChildAt(0, tcPr) }
	p := baseP.Copy()
	for _, ch := range p.ChildElements() { if isTag(ch, "r") { p.RemoveChild(ch) } }
	r := etree.NewElement("w:r"); t := etree.NewElement("w:t"); t.CreateAttr("xml:space", "preserve"); t.SetText(value)
	r.AddChild(t); p.AddChild(r); tc.AddChild(p)
}

func setCellTextMultiPara(tc *etree.Element, rIdx, cIdx int, value string) {
	tcPr := firstDesc(tc, "tcPr"); baseP := firstDesc(tc, "p")
	if baseP == nil { baseP = etree.NewElement("w:p") }
	for _, ch := range tc.ChildElements() { if !isTag(ch, "tcPr") { tc.RemoveChild(ch) } }
	if tcPr == nil { tcPr = etree.NewElement("w:tcPr"); tc.InsertChildAt(0, tcPr) }
	lines := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
	for _, line := range lines {
		p := baseP.Copy(); pPr := firstDesc(p, "pPr")
		if pPr == nil { pPr = etree.NewElement("w:pPr"); p.InsertChildAt(0, pPr) }
		spacing := etree.NewElement("w:spacing"); spacing.CreateAttr("w:after", "0"); spacing.CreateAttr("w:before", "0"); spacing.CreateAttr("w:line", "240"); spacing.CreateAttr("w:lineRule", "auto")
		pPr.AddChild(spacing)
		for _, ch := range p.ChildElements() { if isTag(ch, "r") { p.RemoveChild(ch) } }
		r := etree.NewElement("w:r"); t := etree.NewElement("w:t"); t.CreateAttr("xml:space", "preserve"); t.SetText(line)
		r.AddChild(t); p.AddChild(r); tc.AddChild(p)
	}
}

func setParagraphText(p *etree.Element, text string) {
	t := firstDesc(p, "t")
	if t != nil {
		t.SetText(text)
		for _, otherT := range findDesc(p, "t") { if otherT != t { otherT.SetText("") } }
	} else {
		for _, ch := range p.ChildElements() { if isTag(ch, "r") { p.RemoveChild(ch) } }
		r := etree.NewElement("w:r"); t = etree.NewElement("w:t"); t.SetText(text)
		r.AddChild(t); p.AddChild(r)
	}
}

func getCell(row *etree.Element, idx int) *etree.Element {
	cols := make([]*etree.Element, 0)
	for _, c := range row.ChildElements() { if isTag(c, "tc") { cols = append(cols, c) } }
	if idx >= 0 && idx < len(cols) { return cols[idx] }
	return nil
}

func buildRequestExample(ep endpointEntry, spec *swaggerSpec) string {
	for _, p := range ep.Op.Parameters { if p.In == "body" && p.Schema != nil { return renderSchema(p.Schema, spec, 0) } }
	return "{}"
}

func buildResponseExample(ep endpointEntry, success bool) string {
	if success { return "{\n  \"success\": true,\n  \"message\": \"Data Saved Successfully\"\n}" }
	return "{\n  \"success\": false,\n  \"message\": \"authorization header is required\"\n}"
}

func renderSchema(schema *schemaDef, spec *swaggerSpec, depth int) string {
	if depth > 5 { return "{ \"...\" }" }
	if schema.Ref != "" {
		name := refName(schema.Ref)
		if def, ok := spec.Definitions[name]; ok { return renderSchema(&def, spec, depth+1) }
		return "{ \"payload\": \"" + name + "\" }"
	}
	if strings.ToLower(schema.Type) == "array" && schema.Items != nil {
		itStr := renderProperty(*schema.Items, spec, depth+1)
		return "[\n  " + strings.ReplaceAll(itStr, "\n", "\n  ") + "\n]"
	}
	if len(schema.Properties) == 0 { return "{}" }
	keys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties { keys = append(keys, k) }
	sort.Strings(keys)
	lines := []string{"{"}
	for i, k := range keys {
		v := renderProperty(schema.Properties[k], spec, depth+1)
		sfx := ","; if i == len(keys)-1 { sfx = "" }
		lines = append(lines, fmt.Sprintf("  \"%s\": %s%s", k, v, sfx))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func renderProperty(p propertyDef, spec *swaggerSpec, depth int) string {
	if depth > 8 { return "\"...\"" }
	if p.Ref != "" {
		name := refName(p.Ref)
		if def, ok := spec.Definitions[name]; ok { return renderSchema(&def, spec, depth+1) }
		return "\"value\""
	}
	if strings.ToLower(p.Type) == "array" && p.Items != nil {
		itStr := renderProperty(*p.Items, spec, depth+1)
		return "[\n  " + strings.ReplaceAll(itStr, "\n", "\n  ") + "\n]"
	}
	if strings.ToLower(p.Type) == "object" && len(p.Props) > 0 { return renderSchema(&schemaDef{Properties: p.Props}, spec, depth+1) }
	return sampleValue(p.Type)
}

func sampleValue(typ string) string {
	switch strings.ToLower(typ) {
	case "string": return "\"string\""; case "integer", "number": return "0"; case "boolean": return "true"; case "array": return "[]"; default: return "{}"
	}
}

func refName(ref string) string { p := strings.Split(ref, "/"); return p[len(p)-1] }
func joinPath(base, p string) string { return strings.TrimSuffix(base, "/") + p }
func isTag(el *etree.Element, l string) bool {
	if el == nil { return false }; t := el.Tag; if i := strings.Index(t, ":"); i >= 0 { t = t[i+1:] }; return t == l
}
func findDesc(r *etree.Element, l string) []*etree.Element {
	out := []*etree.Element{}; var w func(*etree.Element); w = func(n *etree.Element) {
		for _, c := range n.ChildElements() { if isTag(c, l) { out = append(out, c) }; w(c) }
	}; w(r); return out
}
func firstDesc(r *etree.Element, l string) *etree.Element {
	for _, c := range r.ChildElements() { if isTag(c, l) { return c }; if res := firstDesc(c, l); res != nil { return res } }; return nil
}
