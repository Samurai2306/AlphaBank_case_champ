package pdftext

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

const maxLLMChars = 28000

var (
	parenString = regexp.MustCompile(`\((?:\\.|[^\\)])*\)`)
	hexString   = regexp.MustCompile(`<([0-9A-Fa-f\s]+)>`)
)

// Result describes extraction quality for demos and LLM ingestion.
type Result struct {
	Text       string   `json:"text"`
	Markdown   string   `json:"markdown"`
	Method     string   `json:"method"` // library | heuristic | merged | empty
	Pages      int      `json:"pages"`
	Chars      int      `json:"chars"`
	Quality    string   `json:"quality"` // good | partial | poor
	Warnings   []string `json:"warnings,omitempty"`
}

// Extract pulls readable text from a PDF (library first, heuristic fallback).
func Extract(raw []byte) string {
	return ExtractDetailed(raw).Text
}

// ExtractDetailed runs a multi-strategy extract and returns LLM-ready markdown.
func ExtractDetailed(raw []byte) Result {
	if len(raw) == 0 {
		return Result{Method: "empty", Quality: "poor", Warnings: []string{"empty file"}}
	}

	libText, pages, libErr := extractWithLibrary(raw)
	heurText := extractHeuristic(raw)

	libText = cleanup(libText)
	heurText = cleanup(heurText)

	var warnings []string
	if libErr != nil {
		warnings = append(warnings, "library: "+libErr.Error())
	}

	method := "empty"
	text := ""
	switch {
	case scoreText(libText) >= scoreText(heurText) && scoreText(libText) >= 20:
		text = libText
		method = "library"
		if scoreText(heurText)*2 > scoreText(libText) && len(heurText) > 0 {
			prefix := heurText
			if len(prefix) > 80 {
				prefix = prefix[:80]
			}
			if !strings.Contains(libText, prefix) {
				text = mergeUnique(libText, heurText)
				method = "merged"
			}
		}
	case scoreText(heurText) >= 20:
		text = heurText
		method = "heuristic"
		warnings = append(warnings, "used heuristic extractor (library weak/empty)")
	default:
		text = pickLonger(libText, heurText)
		if strings.TrimSpace(text) == "" {
			return Result{Method: "empty", Quality: "poor", Pages: pages, Warnings: append(warnings, "no extractable text — likely scanned PDF without OCR")}
		}
		method = "partial"
		warnings = append(warnings, "low text yield")
	}

	text = limitRunes(text, maxLLMChars)
	md := FormatForLLM(text, pages, method)
	chars := len([]rune(text))
	quality := "good"
	if chars < 200 {
		quality = "poor"
	} else if chars < 800 || method == "heuristic" || method == "partial" {
		quality = "partial"
	}

	return Result{
		Text:     text,
		Markdown: md,
		Method:   method,
		Pages:    pages,
		Chars:    chars,
		Quality:  quality,
		Warnings: warnings,
	}
}

func extractWithLibrary(raw []byte) (string, int, error) {
	r, err := pdf.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return "", 0, err
	}
	n := r.NumPage()
	var b strings.Builder

	// Prefer structured row extraction (better word order for contracts).
	rowOK := 0
	for i := 1; i <= n; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		rows, err := p.GetTextByRow()
		if err != nil || len(rows) == 0 {
			continue
		}
		rowOK++
		b.WriteString(fmt.Sprintf("\n\n## Страница %d\n\n", i))
		for _, row := range rows {
			var line []string
			for _, w := range row.Content {
				s := strings.TrimSpace(w.S)
				if s != "" {
					line = append(line, s)
				}
			}
			if len(line) > 0 {
				b.WriteString(strings.Join(line, " "))
				b.WriteByte('\n')
			}
		}
	}
	if rowOK > 0 && len([]rune(b.String())) >= 40 {
		return b.String(), n, nil
	}

	// Fallback: GetPlainText
	plain, err := r.GetPlainText()
	if err != nil {
		return "", n, err
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(plain); err != nil {
		return "", n, err
	}
	return buf.String(), n, nil
}

// FormatForLLM wraps raw extract into a stable markdown envelope for the model.
func FormatForLLM(text string, pages int, method string) string {
	text = strings.TrimSpace(text)
	var b strings.Builder
	b.WriteString("# Извлечённый текст договора\n\n")
	b.WriteString(fmt.Sprintf("- Метод извлечения: `%s`\n", method))
	if pages > 0 {
		b.WriteString(fmt.Sprintf("- Страниц в PDF: %d\n", pages))
	}
	b.WriteString(fmt.Sprintf("- Символов: %d\n\n", len([]rune(text))))
	b.WriteString("---\n\n")
	b.WriteString(text)
	b.WriteString("\n")
	return b.String()
}

func extractHeuristic(raw []byte) string {
	var parts []string
	parts = append(parts, extractParenStrings(raw)...)
	parts = append(parts, extractHexStrings(raw)...)
	parts = append(parts, extractUTF16(raw)...)
	parts = append(parts, extractASCIIRuns(raw)...)
	parts = append(parts, extractCP1251Runs(raw)...)
	return strings.Join(parts, "\n")
}

func extractParenStrings(raw []byte) []string {
	var out []string
	for _, m := range parenString.FindAll(raw, -1) {
		s := decodePDFLiteral(m)
		if useful(s) {
			out = append(out, s)
		}
	}
	return out
}

func extractHexStrings(raw []byte) []string {
	var out []string
	for _, m := range hexString.FindAllStringSubmatch(string(raw), -1) {
		if len(m) < 2 {
			continue
		}
		hex := strings.ReplaceAll(m[1], " ", "")
		if len(hex) < 8 || len(hex)%2 != 0 {
			continue
		}
		buf := make([]byte, 0, len(hex)/2)
		for i := 0; i+1 < len(hex); i += 2 {
			var v byte
			fmt.Sscanf(hex[i:i+2], "%02x", &v)
			buf = append(buf, v)
		}
		// UTF-16BE often starts with BOM FE FF
		if len(buf) >= 2 && buf[0] == 0xFE && buf[1] == 0xFF {
			var u16 []uint16
			for i := 2; i+1 < len(buf); i += 2 {
				u16 = append(u16, uint16(buf[i])<<8|uint16(buf[i+1]))
			}
			s := string(utf16.Decode(u16))
			if useful(s) {
				out = append(out, s)
			}
			continue
		}
		if s := string(buf); utf8.ValidString(s) && useful(s) {
			out = append(out, s)
		}
	}
	return out
}

func decodePDFLiteral(m []byte) string {
	if len(m) < 2 {
		return ""
	}
	inner := m[1 : len(m)-1]
	var b strings.Builder
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' && i+1 < len(inner) {
			i++
			switch inner[i] {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '(', ')', '\\':
				b.WriteByte(inner[i])
			default:
				b.WriteByte(inner[i])
			}
			continue
		}
		b.WriteByte(inner[i])
	}
	return b.String()
}

func extractUTF16(raw []byte) []string {
	var out []string
	for i := 0; i+3 < len(raw); i++ {
		if raw[i] != 0xFE || raw[i+1] != 0xFF {
			continue
		}
		j := i + 2
		var u16 []uint16
		for j+1 < len(raw) && j < i+16000 {
			hi, lo := raw[j], raw[j+1]
			if hi == 0 && lo == 0 {
				break
			}
			u16 = append(u16, uint16(hi)<<8|uint16(lo))
			j += 2
		}
		if len(u16) >= 4 {
			s := string(utf16.Decode(u16))
			if useful(s) {
				out = append(out, s)
			}
		}
		i = j
	}
	return out
}

func extractASCIIRuns(raw []byte) []string {
	var out []string
	var b strings.Builder
	flush := func() {
		s := strings.TrimSpace(b.String())
		b.Reset()
		if useful(s) && looksLikeProse(s) {
			out = append(out, s)
		}
	}
	for _, c := range raw {
		if c >= 32 && c < 127 || c == '\n' || c == '\r' || c == '\t' {
			b.WriteByte(c)
			if b.Len() > 500 {
				flush()
			}
		} else {
			flush()
		}
	}
	flush()
	return out
}

// extractCP1251Runs recovers Windows-1251 Russian text often embedded in older PDFs.
func extractCP1251Runs(raw []byte) []string {
	var out []string
	var buf []byte
	flush := func() {
		if len(buf) < 16 {
			buf = buf[:0]
			return
		}
		s := decodeCP1251(buf)
		buf = buf[:0]
		if useful(s) && looksLikeProse(s) {
			out = append(out, s)
		}
	}
	for _, c := range raw {
		if c >= 0xC0 || (c >= 32 && c < 127) || c == '\n' || c == ' ' {
			buf = append(buf, c)
			if len(buf) > 600 {
				flush()
			}
		} else {
			flush()
		}
	}
	flush()
	return out
}

func decodeCP1251(b []byte) string {
	var r []rune
	for _, c := range b {
		switch {
		case c < 0x80:
			r = append(r, rune(c))
		case c >= 0xC0 && c <= 0xFF:
			// А-Яа-я in CP1251 map linearly onto Unicode Cyrillic block.
			r = append(r, 'А'+rune(int(c)-0xC0))
		case c == 0xA8:
			r = append(r, 'Ё')
		case c == 0xB8:
			r = append(r, 'ё')
		default:
			r = append(r, ' ')
		}
	}
	return string(r)
}

func useful(s string) bool {
	s = strings.TrimSpace(s)
	if len([]rune(s)) < 12 {
		return false
	}
	letters := 0
	cyr := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters++
			if r >= 'А' && r <= 'я' || r == 'ё' || r == 'Ё' {
				cyr++
			}
		}
	}
	return letters >= 8 || cyr >= 6
}

func looksLikeProse(s string) bool {
	lower := strings.ToLower(s)
	noise := []string{"endobj", "endstream", "xref", "/type", "/font", "obj\n", "stream\n", "/length"}
	for _, n := range noise {
		if strings.Contains(lower, n) {
			return false
		}
	}
	if strings.Contains(lower, "аренд") || strings.Contains(lower, "договор") ||
		strings.Contains(lower, "штраф") || strings.Contains(lower, "сторон") ||
		strings.Contains(lower, "плат") || strings.Contains(lower, "кабинет") ||
		strings.Contains(lower, "арендатор") || strings.Contains(lower, "обязатель") {
		return true
	}
	spaces := strings.Count(s, " ")
	return spaces >= 3 && utf8.ValidString(s)
}

func cleanup(s string) string {
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\t", " ")
	lines := strings.Split(s, "\n")
	var keep []string
	seen := map[string]bool{}
	for _, ln := range lines {
		ln = strings.Join(strings.Fields(ln), " ")
		ln = strings.TrimSpace(ln)
		if ln == "" || seen[ln] {
			continue
		}
		seen[ln] = true
		keep = append(keep, ln)
	}
	return strings.Join(keep, "\n")
}

func scoreText(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	score := len([]rune(s))
	lower := strings.ToLower(s)
	for _, k := range []string{"договор", "аренд", "сторон", "плат", "штраф", "обязан", "кабинет", "срок"} {
		if strings.Contains(lower, k) {
			score += 120
		}
	}
	cyr := 0
	for _, r := range s {
		if r >= 'А' && r <= 'я' {
			cyr++
		}
	}
	score += cyr / 2
	return score
}

func mergeUnique(a, b string) string {
	seen := map[string]bool{}
	var out []string
	for _, block := range []string{a, b} {
		for _, ln := range strings.Split(block, "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" || seen[ln] {
				continue
			}
			seen[ln] = true
			out = append(out, ln)
		}
	}
	return strings.Join(out, "\n")
}

func pickLonger(a, b string) string {
	if len([]rune(a)) >= len([]rune(b)) {
		return a
	}
	return b
}

func limitRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "\n…[текст обрезан для контекста модели]"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HasText reports whether extraction found meaningful content.
func HasText(raw []byte) bool {
	return ExtractDetailed(raw).Chars >= 40
}

// SampleContractBytes builds a tiny PDF-like payload for tests.
func SampleContractBytes() []byte {
	body := `%PDF-1.4
1 0 obj<<>>endobj
2 0 obj<< /Length 200 >>stream
BT
(Договор аренды кабинета. Арендодатель вправе в одностороннем порядке изменить стоимость аренды.) Tj
(При расторжении по инициативе арендатора взыскивается штраф за весь срок.) Tj
(Споры рассматриваются по месту нахождения арендодателя.) Tj
ET
endstream
endobj
trailer<<>>
%%EOF`
	return bytes.ReplaceAll([]byte(body), []byte("\r\n"), []byte("\n"))
}
