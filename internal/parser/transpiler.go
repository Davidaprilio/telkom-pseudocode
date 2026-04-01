package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type block struct {
	kind string
	line int
}

func Transpile(source string) (string, error) {
	lines := strings.Split(source, "\n")
	decls := make([]string, 0)
	body := make([]string, 0)
	stack := make([]block, 0)

	section := "none"
	indent := 1
	loopCounter := 0
	indentKind := ""
	indentUnit := 0

	for i, raw := range lines {
		lineNo := i + 1
		cleanRaw := stripCurlyComment(raw)
		line := strings.TrimSpace(cleanRaw)
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "program ") {
			continue
		}

		switch lower {
		case "kamus":
			section = "kamus"
			continue
		case "algoritma":
			section = "algoritma"
			continue
		case "endalgoritma":
			section = "done"
			continue
		}

		if section == "kamus" {
			indentLevel, kind, unit, err := computeIndentLevel(cleanRaw, lineNo, indentKind, indentUnit)
			if err != nil {
				return "", err
			}
			indentKind, indentUnit = kind, unit
			if indentLevel != 1 {
				return "", fmt.Errorf("baris %d: deklarasi di blok kamus harus indent level 1", lineNo)
			}

			decl, err := parseDeclaration(line, lineNo)
			if err != nil {
				return "", err
			}
			decls = append(decls, indentLine(1, decl))
			continue
		}

		if section != "algoritma" {
			return "", fmt.Errorf("baris %d: konten harus berada di bawah blok kamus/algoritma", lineNo)
		}

		indentLevel, kind, unit, err := computeIndentLevel(cleanRaw, lineNo, indentKind, indentUnit)
		if err != nil {
			return "", err
		}
		indentKind, indentUnit = kind, unit

		expectedLevel := len(stack)
		if (strings.HasPrefix(lower, "else if ") && strings.HasSuffix(lower, " then")) || lower == "else" || lower == "endif" || lower == "endfor" || lower == "akhir-ulangi" {
			expectedLevel = len(stack) - 1
		}
		if expectedLevel < 0 {
			expectedLevel = 0
		}
		if indentLevel != expectedLevel {
			return "", fmt.Errorf("baris %d: indentasi tidak konsisten, expected level %d tetapi mendapat level %d", lineNo, expectedLevel, indentLevel)
		}

		if strings.HasPrefix(lower, "if ") && strings.HasSuffix(lower, " then") {
			cond := strings.TrimSpace(line[3 : len(line)-5])
			cond = normalizeCondition(cond)
			body = append(body, indentLine(indent, fmt.Sprintf("if %s {", cond)))
			stack = append(stack, block{kind: "if", line: lineNo})
			indent++
			continue
		}

		if strings.HasPrefix(lower, "else if ") && strings.HasSuffix(lower, " then") {
			if len(stack) == 0 || stack[len(stack)-1].kind != "if" {
				return "", fmt.Errorf("baris %d: else if tanpa blok if", lineNo)
			}
			cond := strings.TrimSpace(line[len("else if ") : len(line)-len(" then")])
			if cond == "" {
				return "", fmt.Errorf("baris %d: kondisi else if kosong", lineNo)
			}
			indent--
			body = append(body, indentLine(indent, fmt.Sprintf("} else if %s {", normalizeCondition(cond))))
			indent++
			continue
		}

		if lower == "else" {
			if len(stack) == 0 || stack[len(stack)-1].kind != "if" {
				return "", fmt.Errorf("baris %d: else tanpa blok if", lineNo)
			}
			indent--
			body = append(body, indentLine(indent, "} else {"))
			indent++
			continue
		}

		if lower == "endif" {
			if len(stack) == 0 || stack[len(stack)-1].kind != "if" {
				return "", fmt.Errorf("baris %d: endif tanpa if", lineNo)
			}
			indent--
			body = append(body, indentLine(indent, "}"))
			stack = stack[:len(stack)-1]
			continue
		}

		if strings.HasPrefix(lower, "for ") && strings.HasSuffix(lower, " do") {
			loopLine, err := parseForToLoop(line, lineNo)
			if err != nil {
				return "", err
			}
			body = append(body, indentLine(indent, loopLine))
			stack = append(stack, block{kind: "loop", line: lineNo})
			indent++
			continue
		}

		if strings.HasPrefix(lower, "ulangi ") && strings.HasSuffix(lower, " kali") {
			expr := strings.TrimSpace(line[7 : len(line)-5])
			if expr == "" {
				return "", fmt.Errorf("baris %d: jumlah perulangan kosong", lineNo)
			}
			loopCounter++
			idxName := fmt.Sprintf("i%d", loopCounter)
			body = append(body, indentLine(indent, fmt.Sprintf("for %s := 0; %s < (%s); %s++ {", idxName, idxName, normalizeExpression(expr), idxName)))
			stack = append(stack, block{kind: "loop", line: lineNo})
			indent++
			continue
		}

		if lower == "akhir-ulangi" || lower == "endfor" {
			if len(stack) == 0 || stack[len(stack)-1].kind != "loop" {
				return "", fmt.Errorf("baris %d: akhir-ulangi tanpa ulangi", lineNo)
			}
			indent--
			body = append(body, indentLine(indent, "}"))
			stack = stack[:len(stack)-1]
			continue
		}

		if (strings.HasPrefix(lower, "tulis(") || strings.HasPrefix(lower, "output(")) && strings.HasSuffix(line, ")") {
			fnLen := len("tulis(")
			if strings.HasPrefix(lower, "output(") {
				fnLen = len("output(")
			}
			inner := strings.TrimSpace(line[fnLen : len(line)-1])
			body = append(body, indentLine(indent, fmt.Sprintf("fmt.Println(%s)", normalizeExpression(inner))))
			continue
		}

		if (strings.HasPrefix(lower, "baca(") || strings.HasPrefix(lower, "input(")) && strings.HasSuffix(line, ")") {
			fnLen := len("baca(")
			if strings.HasPrefix(lower, "input(") {
				fnLen = len("input(")
			}
			inner := strings.TrimSpace(line[fnLen : len(line)-1])
			scanArgs, err := parseInputArgs(inner, lineNo)
			if err != nil {
				return "", err
			}
			body = append(body, indentLine(indent, fmt.Sprintf("fmt.Scanln(%s)", scanArgs)))
			continue
		}

		if strings.Contains(line, "<-") {
			parts := strings.SplitN(line, "<-", 2)
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			if left == "" || right == "" {
				return "", fmt.Errorf("baris %d: assignment tidak valid", lineNo)
			}
			body = append(body, indentLine(indent, fmt.Sprintf("%s = %s", left, normalizeExpression(right))))
			continue
		}

		return "", fmt.Errorf("baris %d: sintaks tidak dikenali: %s", lineNo, line)
	}

	if len(stack) > 0 {
		open := stack[len(stack)-1]
		switch open.kind {
		case "if":
			return "", fmt.Errorf("baris %d: blok jika belum ditutup (butuh akhir-jika)", open.line)
		case "loop":
			return "", fmt.Errorf("baris %d: blok ulangi belum ditutup (butuh akhir-ulangi)", open.line)
		default:
			return "", fmt.Errorf("baris %d: blok belum ditutup", open.line)
		}
	}

	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import \"fmt\"\n\n")
	b.WriteString("func main() {\n")

	for _, d := range decls {
		b.WriteString(d)
		b.WriteString("\n")
	}
	if len(decls) > 0 && len(body) > 0 {
		b.WriteString("\n")
	}
	for _, s := range body {
		b.WriteString(s)
		b.WriteString("\n")
	}

	b.WriteString("}\n")
	return b.String(), nil
}

func parseDeclaration(line string, lineNo int) (string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: deklarasi kamus harus 'nama: tipe'", lineNo)
	}
	namesPart := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(strings.ToLower(parts[1]))
	if namesPart == "" {
		return "", fmt.Errorf("baris %d: nama variabel kosong", lineNo)
	}
	goType, ok := mapType(typ)
	if !ok {
		return "", fmt.Errorf("baris %d: tipe data tidak dikenal '%s'", lineNo, typ)
	}
	rawNames := strings.Split(namesPart, ",")
	names := make([]string, 0, len(rawNames))
	for _, raw := range rawNames {
		name := strings.TrimSpace(raw)
		if name == "" {
			return "", fmt.Errorf("baris %d: daftar nama variabel tidak valid", lineNo)
		}
		names = append(names, name)
	}
	return fmt.Sprintf("var %s %s", strings.Join(names, ", "), goType), nil
}

func mapType(pseudoType string) (string, bool) {
	switch pseudoType {
	case "integer":
		return "int", true
	case "real":
		return "float64", true
	case "string":
		return "string", true
	case "boolean":
		return "bool", true
	default:
		return "", false
	}
}

func indentLine(level int, content string) string {
	return strings.Repeat("\t", level) + content
}

func normalizeExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	expr = strings.ReplaceAll(expr, "'", "\"")
	return expr
}

func normalizeCondition(cond string) string {
	cond = normalizeExpression(cond)
	re := regexp.MustCompile(`(^|[^<>!])=([^=])`)
	cond = re.ReplaceAllString(cond, `${1}==${2}`)
	return cond
}

func stripCurlyComment(line string) string {
	commentRe := regexp.MustCompile(`\{[^{}]*\}`)
	for commentRe.MatchString(line) {
		line = commentRe.ReplaceAllString(line, "")
	}
	return line
}

func parseInputArgs(inner string, lineNo int) (string, error) {
	if inner == "" {
		return "", fmt.Errorf("baris %d: argumen input/baca kosong", lineNo)
	}
	rawArgs := strings.Split(inner, ",")
	args := make([]string, 0, len(rawArgs))
	for _, raw := range rawArgs {
		name := strings.TrimSpace(raw)
		if name == "" {
			return "", fmt.Errorf("baris %d: daftar variabel input/baca tidak valid", lineNo)
		}
		args = append(args, "&"+name)
	}
	return strings.Join(args, ", "), nil
}

func parseForToLoop(line string, lineNo int) (string, error) {
	inner := strings.TrimSpace(line[len("for ") : len(line)-len(" do")])
	if inner == "" {
		return "", fmt.Errorf("baris %d: sintaks for tidak valid", lineNo)
	}
	parts := strings.SplitN(inner, " to ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: gunakan format 'for <var> to <batas> do'", lineNo)
	}
	iterVar := strings.TrimSpace(parts[0])
	limit := strings.TrimSpace(parts[1])
	if iterVar == "" || limit == "" {
		return "", fmt.Errorf("baris %d: variabel atau batas for kosong", lineNo)
	}
	return fmt.Sprintf("for %s <= (%s) {", iterVar, normalizeExpression(limit)), nil
}

func computeIndentLevel(rawLine string, lineNo int, currentKind string, currentUnit int) (int, string, int, error) {
	leading := leadingWhitespace(rawLine)
	if strings.Contains(leading, " ") && strings.Contains(leading, "\t") {
		return 0, currentKind, currentUnit, fmt.Errorf("baris %d: indentasi tidak boleh mencampur spasi dan tab", lineNo)
	}

	if leading == "" {
		return 0, currentKind, currentUnit, nil
	}

	kind := "space"
	if strings.Contains(leading, "\t") {
		kind = "tab"
	}

	if currentKind != "" && kind != currentKind {
		return 0, currentKind, currentUnit, fmt.Errorf("baris %d: indentasi harus konsisten, sebelumnya %s sekarang %s", lineNo, currentKind, kind)
	}

	width := len(leading)
	unit := currentUnit
	if kind == "tab" {
		unit = 1
	}
	if unit == 0 {
		unit = width
	}

	if width%unit != 0 {
		return 0, kind, unit, fmt.Errorf("baris %d: indentasi tidak konsisten dengan unit %d", lineNo, unit)
	}

	return width / unit, kind, unit, nil
}

func leadingWhitespace(s string) string {
	idx := 0
	for idx < len(s) {
		if s[idx] != ' ' && s[idx] != '\t' {
			break
		}
		idx++
	}
	return s[:idx]
}
