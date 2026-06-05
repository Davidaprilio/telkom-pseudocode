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

type routine struct {
	name       string
	kind       string
	header     string
	decls      []string
	body       []string
	stack      []block
	indent     int
	line       int
	hasAlgo    bool
	closed     bool
	isMain     bool
	loopNumber int
}

type transpiler struct {
	global     []string
	routines   []routine
	current    *routine
	section    string
	indentKind string
	indentUnit int
	record     *recordBuilder
}

type recordBuilder struct {
	name   string
	fields []string
	line   int
}

func Transpile(source string) (string, error) {
	t := &transpiler{
		global:   make([]string, 0),
		routines: make([]routine, 0),
		section:  "top",
	}

	lines := strings.Split(source, "\n")
	for i, raw := range lines {
		lineNo := i + 1
		cleanRaw := stripCurlyComment(raw)
		line := strings.TrimSpace(cleanRaw)
		if line == "" {
			continue
		}
		if err := t.handleLine(cleanRaw, line, lineNo); err != nil {
			return "", err
		}
	}

	if t.record != nil {
		return "", fmt.Errorf("baris %d: type record belum ditutup dengan >", t.record.line)
	}
	if t.current != nil {
		if err := t.finishCurrent(); err != nil {
			return "", err
		}
	}

	return t.emit()
}

func (t *transpiler) handleLine(raw, line string, lineNo int) error {
	lower := strings.ToLower(line)

	if t.record != nil {
		if line == ">" {
			t.global = append(t.global, fmt.Sprintf("type %s struct {", t.record.name))
			t.global = append(t.global, t.record.fields...)
			t.global = append(t.global, "}")
			t.record = nil
			return nil
		}
		field, err := parseRecordField(line, lineNo)
		if err != nil {
			return err
		}
		t.record.fields = append(t.record.fields, indentLine(1, field))
		return nil
	}

	if strings.HasPrefix(lower, "type ") && strings.HasSuffix(line, "<") {
		if t.current != nil {
			return fmt.Errorf("baris %d: type record hanya didukung di kamus global/top-level", lineNo)
		}
		name := strings.TrimSpace(line[len("type ") : len(line)-1])
		if name == "" {
			return fmt.Errorf("baris %d: nama type kosong", lineNo)
		}
		t.record = &recordBuilder{name: name, fields: make([]string, 0), line: lineNo}
		return nil
	}

	if strings.HasPrefix(lower, "constant ") {
		constLine, err := parseConstant(line, lineNo)
		if err != nil {
			return err
		}
		t.addDeclaration(constLine)
		return nil
	}

	if strings.HasPrefix(lower, "type ") {
		typeLine, err := parseTypeAlias(line, lineNo)
		if err != nil {
			return err
		}
		t.addDeclaration(typeLine)
		return nil
	}

	if lower == "kamus global" {
		if t.current != nil {
			return fmt.Errorf("baris %d: kamus global tidak boleh berada di dalam blok", lineNo)
		}
		t.section = "global"
		return nil
	}

	if strings.HasPrefix(lower, "program ") {
		if t.current != nil {
			return fmt.Errorf("baris %d: blok sebelumnya belum ditutup", lineNo)
		}
		name := strings.TrimSpace(line[len("program "):])
		if name == "" {
			return fmt.Errorf("baris %d: nama program kosong", lineNo)
		}
		t.current = &routine{name: name, kind: "program", header: "func main() {", decls: make([]string, 0), body: make([]string, 0), indent: 1, line: lineNo, isMain: true}
		t.section = "routine"
		return nil
	}

	if strings.HasPrefix(lower, "procedure ") {
		if t.current != nil {
			return fmt.Errorf("baris %d: blok sebelumnya belum ditutup", lineNo)
		}
		header, name, err := parseProcedureHeader(line, lineNo)
		if err != nil {
			return err
		}
		t.current = &routine{name: name, kind: "procedure", header: header, decls: make([]string, 0), body: make([]string, 0), indent: 1, line: lineNo}
		t.section = "routine"
		return nil
	}

	if strings.HasPrefix(lower, "function ") {
		if t.current != nil {
			return fmt.Errorf("baris %d: blok sebelumnya belum ditutup", lineNo)
		}
		header, name, err := parseFunctionHeader(line, lineNo)
		if err != nil {
			return err
		}
		t.current = &routine{name: name, kind: "function", header: header, decls: make([]string, 0), body: make([]string, 0), indent: 1, line: lineNo}
		t.section = "routine"
		return nil
	}

	if t.current == nil {
		if t.section == "global" {
			decl, err := parseDeclaration(line, lineNo)
			if err != nil {
				return err
			}
			t.global = append(t.global, decl)
			return nil
		}
		return fmt.Errorf("baris %d: konten harus berada di dalam program/procedure/function atau kamus global", lineNo)
	}

	switch lower {
	case "kamus":
		t.section = "kamus"
		return nil
	case "algoritma":
		t.section = "algoritma"
		t.current.hasAlgo = true
		return nil
	case "endalgoritma":
		t.section = "routine"
		return nil
	case "endprogram", "endprocedure", "endfunction":
		return t.closeRoutine(lineNo, lower)
	}

	if t.section == "kamus" {
		indentLevel, err := t.indentLevel(raw, lineNo)
		if err != nil {
			return err
		}
		if indentLevel != 1 {
			return fmt.Errorf("baris %d: deklarasi di blok kamus harus indent level 1", lineNo)
		}
		decl, err := parseDeclaration(line, lineNo)
		if err != nil {
			return err
		}
		t.current.decls = append(t.current.decls, indentLine(1, decl))
		return nil
	}

	if t.section != "algoritma" {
		return fmt.Errorf("baris %d: konten harus berada di bawah blok kamus/algoritma", lineNo)
	}

	return t.addStatement(raw, line, lineNo)
}

func (t *transpiler) addDeclaration(decl string) {
	if t.current == nil {
		t.global = append(t.global, decl)
		return
	}
	t.current.decls = append(t.current.decls, indentLine(1, decl))
}

func (t *transpiler) closeRoutine(lineNo int, closer string) error {
	if t.current == nil {
		return fmt.Errorf("baris %d: %s tanpa pembuka", lineNo, closer)
	}
	expectedCloser := "end" + t.current.kind
	if closer != expectedCloser {
		return fmt.Errorf("baris %d: penutup %s tidak sesuai", lineNo, closer)
	}
	return t.finishCurrent()
}

func (t *transpiler) finishCurrent() error {
	if len(t.current.stack) > 0 {
		open := t.current.stack[len(t.current.stack)-1]
		return fmt.Errorf("baris %d: blok %s belum ditutup", open.line, open.kind)
	}
	t.current.closed = true
	t.routines = append(t.routines, *t.current)
	t.current = nil
	t.section = "top"
	return nil
}

func (t *transpiler) addStatement(raw, line string, lineNo int) error {
	indentLevel, err := t.indentLevel(raw, lineNo)
	if err != nil {
		return err
	}

	lower := strings.ToLower(line)
	expectedLevel := len(t.current.stack) + 1
	if isClosingStatement(lower) {
		expectedLevel = len(t.current.stack)
	}
	if expectedLevel < 1 {
		expectedLevel = 1
	}
	if indentLevel != expectedLevel {
		return fmt.Errorf("baris %d: indentasi tidak konsisten, expected level %d tetapi mendapat level %d", lineNo, expectedLevel, indentLevel)
	}

	indent := t.current.indent

	if strings.HasPrefix(lower, "if ") && strings.HasSuffix(lower, " then") {
		cond := strings.TrimSpace(line[3 : len(line)-5])
		if cond == "" {
			return fmt.Errorf("baris %d: kondisi if kosong", lineNo)
		}
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("if %s {", normalizeCondition(cond))))
		t.current.stack = append(t.current.stack, block{kind: "if", line: lineNo})
		t.current.indent++
		return nil
	}

	if strings.HasPrefix(lower, "else if ") && strings.HasSuffix(lower, " then") {
		if len(t.current.stack) == 0 || t.current.stack[len(t.current.stack)-1].kind != "if" {
			return fmt.Errorf("baris %d: else if tanpa blok if", lineNo)
		}
		cond := strings.TrimSpace(line[len("else if ") : len(line)-len(" then")])
		t.current.indent--
		t.current.body = append(t.current.body, indentLine(t.current.indent, fmt.Sprintf("} else if %s {", normalizeCondition(cond))))
		t.current.indent++
		return nil
	}

	if lower == "else" {
		if len(t.current.stack) == 0 || t.current.stack[len(t.current.stack)-1].kind != "if" {
			return fmt.Errorf("baris %d: else tanpa blok if", lineNo)
		}
		t.current.indent--
		t.current.body = append(t.current.body, indentLine(t.current.indent, "} else {"))
		t.current.indent++
		return nil
	}

	if lower == "endif" {
		if err := t.closeStatementBlock("if", lineNo, "endif tanpa if"); err != nil {
			return err
		}
		return nil
	}

	if strings.HasPrefix(lower, "while ") && strings.HasSuffix(lower, " do") {
		cond := strings.TrimSpace(line[len("while ") : len(line)-len(" do")])
		if cond == "" {
			return fmt.Errorf("baris %d: kondisi while kosong", lineNo)
		}
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("for %s {", normalizeCondition(cond))))
		t.current.stack = append(t.current.stack, block{kind: "while", line: lineNo})
		t.current.indent++
		return nil
	}

	if lower == "endwhile" {
		if err := t.closeStatementBlock("while", lineNo, "endwhile tanpa while"); err != nil {
			return err
		}
		return nil
	}

	if strings.HasPrefix(lower, "for ") && strings.HasSuffix(lower, " do") {
		loopLine, err := parseForToLoop(line, lineNo)
		if err != nil {
			return err
		}
		t.current.body = append(t.current.body, indentLine(indent, loopLine))
		t.current.stack = append(t.current.stack, block{kind: "for", line: lineNo})
		t.current.indent++
		return nil
	}

	if strings.HasPrefix(lower, "ulangi ") && strings.HasSuffix(lower, " kali") {
		expr := strings.TrimSpace(line[7 : len(line)-5])
		if expr == "" {
			return fmt.Errorf("baris %d: jumlah perulangan kosong", lineNo)
		}
		t.current.loopNumber++
		idxName := fmt.Sprintf("i%d", t.current.loopNumber)
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("for %s := 0; %s < (%s); %s++ {", idxName, idxName, normalizeExpression(expr), idxName)))
		t.current.stack = append(t.current.stack, block{kind: "repeat", line: lineNo})
		t.current.indent++
		return nil
	}

	if lower == "endfor" || lower == "akhir-ulangi" {
		if len(t.current.stack) == 0 || (t.current.stack[len(t.current.stack)-1].kind != "for" && t.current.stack[len(t.current.stack)-1].kind != "repeat") {
			return fmt.Errorf("baris %d: %s tanpa pembuka perulangan", lineNo, lower)
		}
		t.current.indent--
		t.current.body = append(t.current.body, indentLine(t.current.indent, "}"))
		t.current.stack = t.current.stack[:len(t.current.stack)-1]
		return nil
	}

	if (strings.HasPrefix(lower, "tulis(") || strings.HasPrefix(lower, "output(")) && strings.HasSuffix(line, ")") {
		fnLen := len("tulis(")
		if strings.HasPrefix(lower, "output(") {
			fnLen = len("output(")
		}
		inner := strings.TrimSpace(line[fnLen : len(line)-1])
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("fmt.Println(%s)", normalizeExpression(inner))))
		return nil
	}

	if (strings.HasPrefix(lower, "baca(") || strings.HasPrefix(lower, "input(")) && strings.HasSuffix(line, ")") {
		fnLen := len("baca(")
		if strings.HasPrefix(lower, "input(") {
			fnLen = len("input(")
		}
		inner := strings.TrimSpace(line[fnLen : len(line)-1])
		scanArgs, err := parseInputArgs(inner, lineNo)
		if err != nil {
			return err
		}
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("fmt.Scanln(%s)", scanArgs)))
		return nil
	}

	if strings.HasPrefix(lower, "return") {
		expr := strings.TrimSpace(line[len("return"):])
		if expr == "" {
			t.current.body = append(t.current.body, indentLine(indent, "return"))
		} else {
			t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("return %s", normalizeExpression(expr))))
		}
		return nil
	}

	if lower == "break" || lower == "continue" {
		t.current.body = append(t.current.body, indentLine(indent, lower))
		return nil
	}

	if left, right, ok := splitAssignment(line); ok {
		t.current.body = append(t.current.body, indentLine(indent, fmt.Sprintf("%s = %s", left, normalizeExpression(right))))
		return nil
	}

	if isCall(line) {
		t.current.body = append(t.current.body, indentLine(indent, normalizeExpression(line)))
		return nil
	}

	return fmt.Errorf("baris %d: sintaks tidak dikenali: %s", lineNo, line)
}

func (t *transpiler) closeStatementBlock(kind string, lineNo int, message string) error {
	if len(t.current.stack) == 0 || t.current.stack[len(t.current.stack)-1].kind != kind {
		return fmt.Errorf("baris %d: %s", lineNo, message)
	}
	t.current.indent--
	t.current.body = append(t.current.body, indentLine(t.current.indent, "}"))
	t.current.stack = t.current.stack[:len(t.current.stack)-1]
	return nil
}

func (t *transpiler) indentLevel(raw string, lineNo int) (int, error) {
	level, kind, unit, err := computeIndentLevel(raw, lineNo, t.indentKind, t.indentUnit)
	if err != nil {
		return 0, err
	}
	t.indentKind, t.indentUnit = kind, unit
	return level, nil
}

func (t *transpiler) emit() (string, error) {
	hasMain := false
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import \"fmt\"\n\n")

	for _, g := range t.global {
		b.WriteString(g)
		b.WriteString("\n")
	}
	if len(t.global) > 0 {
		b.WriteString("\n")
	}

	for idx, r := range t.routines {
		if r.isMain {
			hasMain = true
		}
		if idx > 0 {
			b.WriteString("\n")
		}
		b.WriteString(r.header)
		b.WriteString("\n")
		for _, d := range r.decls {
			b.WriteString(d)
			b.WriteString("\n")
		}
		if len(r.decls) > 0 && len(r.body) > 0 {
			b.WriteString("\n")
		}
		for _, s := range r.body {
			b.WriteString(s)
			b.WriteString("\n")
		}
		b.WriteString("}\n")
	}

	if !hasMain {
		return "", fmt.Errorf("program utama belum ditemukan")
	}
	return b.String(), nil
}

func parseDeclaration(line string, lineNo int) (string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: deklarasi kamus harus 'nama: tipe'", lineNo)
	}
	namesPart := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])
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

func parseConstant(line string, lineNo int) (string, error) {
	inner := strings.TrimSpace(line[len("constant "):])
	parts := strings.SplitN(inner, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: constant harus 'constant NAMA : tipe = nilai'", lineNo)
	}
	left := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	nameType := strings.Split(left, ":")
	if len(nameType) != 2 {
		return "", fmt.Errorf("baris %d: constant harus menyertakan tipe", lineNo)
	}
	name := strings.TrimSpace(nameType[0])
	typ := strings.TrimSpace(nameType[1])
	goType, ok := mapType(typ)
	if !ok {
		return "", fmt.Errorf("baris %d: tipe data tidak dikenal '%s'", lineNo, typ)
	}
	return fmt.Sprintf("const %s %s = %s", name, goType, normalizeExpression(value)), nil
}

func parseTypeAlias(line string, lineNo int) (string, error) {
	inner := strings.TrimSpace(line[len("type "):])
	parts := strings.SplitN(inner, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: type harus 'type Nama : tipe'", lineNo)
	}
	name := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])
	goType, ok := mapType(typ)
	if !ok {
		return "", fmt.Errorf("baris %d: tipe data tidak dikenal '%s'", lineNo, typ)
	}
	return fmt.Sprintf("type %s %s", name, goType), nil
}

func parseRecordField(line string, lineNo int) (string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: field record harus 'Nama: tipe'", lineNo)
	}
	name := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])
	goType, ok := mapType(typ)
	if !ok {
		return "", fmt.Errorf("baris %d: tipe data tidak dikenal '%s'", lineNo, typ)
	}
	return fmt.Sprintf("%s %s", name, goType), nil
}

func parseProcedureHeader(line string, lineNo int) (string, string, error) {
	name, params, err := parseRoutineSignature(strings.TrimSpace(line[len("procedure "):]), lineNo)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("func %s(%s) {", name, strings.Join(params, ", ")), name, nil
}

func parseFunctionHeader(line string, lineNo int) (string, string, error) {
	inner := strings.TrimSpace(line[len("function "):])
	parts := strings.SplitN(inner, "->", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("baris %d: function harus memiliki tipe return dengan '->'", lineNo)
	}
	name, params, err := parseRoutineSignature(strings.TrimSpace(parts[0]), lineNo)
	if err != nil {
		return "", "", err
	}
	retType := strings.TrimSpace(parts[1])
	goRetType, ok := mapType(retType)
	if !ok {
		return "", "", fmt.Errorf("baris %d: tipe return tidak dikenal '%s'", lineNo, retType)
	}
	return fmt.Sprintf("func %s(%s) %s {", name, strings.Join(params, ", "), goRetType), name, nil
}

func parseRoutineSignature(sig string, lineNo int) (string, []string, error) {
	open := strings.Index(sig, "(")
	close := strings.LastIndex(sig, ")")
	if open <= 0 || close < open {
		return "", nil, fmt.Errorf("baris %d: signature routine tidak valid", lineNo)
	}
	name := strings.TrimSpace(sig[:open])
	paramsText := strings.TrimSpace(sig[open+1 : close])
	params := make([]string, 0)
	if paramsText == "" {
		return name, params, nil
	}
	for _, raw := range splitCSV(paramsText) {
		param, err := parseParameter(raw, lineNo)
		if err != nil {
			return "", nil, err
		}
		params = append(params, param)
	}
	return name, params, nil
}

func parseParameter(raw string, lineNo int) (string, error) {
	text := strings.TrimSpace(raw)
	for _, mode := range []string{"in/out ", "in ", "out "} {
		if strings.HasPrefix(strings.ToLower(text), mode) {
			text = strings.TrimSpace(text[len(mode):])
			break
		}
	}
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("baris %d: parameter harus 'nama : tipe'", lineNo)
	}
	name := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])
	goType, ok := mapType(typ)
	if !ok {
		return "", fmt.Errorf("baris %d: tipe parameter tidak dikenal '%s'", lineNo, typ)
	}
	return fmt.Sprintf("%s %s", name, goType), nil
}

func mapType(pseudoType string) (string, bool) {
	pseudoType = strings.TrimSpace(pseudoType)
	lower := strings.ToLower(pseudoType)
	switch lower {
	case "integer":
		return "int", true
	case "real":
		return "float64", true
	case "string":
		return "string", true
	case "boolean":
		return "bool", true
	case "character", "char":
		return "rune", true
	}
	if strings.HasPrefix(lower, "array") {
		return mapArrayType(pseudoType)
	}
	if isIdentifier(pseudoType) {
		return pseudoType, true
	}
	return "", false
}

func mapArrayType(pseudoType string) (string, bool) {
	reRange := regexp.MustCompile(`(?i)^array\s*\[(.+)\.\.(.+)\]\s+of\s+(.+)$`)
	if m := reRange.FindStringSubmatch(pseudoType); m != nil {
		start := strings.TrimSpace(m[1])
		end := strings.TrimSpace(m[2])
		elem, ok := mapType(strings.TrimSpace(m[3]))
		if !ok {
			return "", false
		}
		if start == "0" {
			return fmt.Sprintf("[%s + 1]%s", normalizeExpression(end), elem), true
		}
		if start == "1" {
			return fmt.Sprintf("[%s + 1]%s", normalizeExpression(end), elem), true
		}
		return fmt.Sprintf("[%s - %s + 1]%s", normalizeExpression(end), normalizeExpression(start), elem), true
	}
	reOpen := regexp.MustCompile(`(?i)^array\s+of\s+(.+)$`)
	if m := reOpen.FindStringSubmatch(pseudoType); m != nil {
		elem, ok := mapType(strings.TrimSpace(m[1]))
		if !ok {
			return "", false
		}
		return "[]" + elem, true
	}
	return "", false
}

func indentLine(level int, content string) string {
	return strings.Repeat("\t", level) + content
}

func normalizeExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	return expr
}

func normalizeCondition(cond string) string {
	cond = normalizeExpression(cond)
	re := regexp.MustCompile(`(^|[^<>!=])=([^=])`)
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
	rawArgs := splitCSV(inner)
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
	return fmt.Sprintf("for ; %s <= (%s); %s++ {", iterVar, normalizeExpression(limit), iterVar), nil
}

func splitAssignment(line string) (string, string, bool) {
	if strings.Contains(line, "<-") {
		parts := strings.SplitN(line, "<-", 2)
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])
		return left, right, left != "" && right != ""
	}
	for i := 0; i < len(line); i++ {
		if line[i] != '=' {
			continue
		}
		prev := byte(0)
		next := byte(0)
		if i > 0 {
			prev = line[i-1]
		}
		if i+1 < len(line) {
			next = line[i+1]
		}
		if prev == '<' || prev == '>' || prev == '!' || prev == '=' || next == '=' {
			continue
		}
		left := strings.TrimSpace(line[:i])
		right := strings.TrimSpace(line[i+1:])
		return left, right, left != "" && right != ""
	}
	return "", "", false
}

func isClosingStatement(lower string) bool {
	return (strings.HasPrefix(lower, "else if ") && strings.HasSuffix(lower, " then")) ||
		lower == "else" ||
		lower == "endif" ||
		lower == "endfor" ||
		lower == "endwhile" ||
		lower == "akhir-ulangi"
}

func isCall(line string) bool {
	open := strings.Index(line, "(")
	return open > 0 && strings.HasSuffix(line, ")") && isIdentifier(strings.TrimSpace(line[:open]))
}

func isIdentifier(name string) bool {
	re := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	return re.MatchString(name)
}

func splitCSV(text string) []string {
	parts := make([]string, 0)
	start := 0
	depth := 0
	inString := false
	var quote byte
	for i := 0; i < len(text); i++ {
		ch := text[i]
		if inString {
			if ch == quote {
				inString = false
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			inString = true
			quote = ch
			continue
		}
		switch ch {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, text[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, text[start:])
	return parts
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
