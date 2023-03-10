package lang

import (
	"go/ast"
	"go/printer"
	"go/token"
	"regexp"
	"strings"
	"unicode"
)

func printNode(node ast.Node, fs *token.FileSet) (string, error) {
	cfg := printer.Config{
		Mode:     printer.UseSpaces,
		Tabwidth: 4,
	}

	var out strings.Builder
	if err := cfg.Fprint(&out, fs, node); err != nil {
		return "", err
	}

	return out.String(), nil
}

func runeIsUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

const lowerToUpper = 'a' - 'A'

func runeToUpper(r rune) rune {
	return r - lowerToUpper
}

func splitCamel(text string) string {
	var builder strings.Builder
	var previousRune rune
	var wordLength int
	for i, r := range text {
		if i == 0 {
			previousRune = runeToUpper(r)
			continue
		}

		switch {
		case runeIsUpper(previousRune) && !runeIsUpper(r) && wordLength > 0:
			// If we have a capital followed by a lower, that capital should
			// begin a word. Throw a space before the runes if there is a word
			// there.
			builder.WriteRune(' ')
			builder.WriteRune(previousRune)
			wordLength = 1
		case !runeIsUpper(previousRune) && runeIsUpper(r):
			// If we have a lower followed by a capital, the capital should
			// begin a word. Throw a space in between the runes. We don't have
			// to check word length because we're writing the previous rune to
			// the previous word, automaticall giving it a length of 1.
			builder.WriteRune(previousRune)
			builder.WriteRune(' ')
			wordLength = 0
		default:
			// Otherwise, just throw the rune onto the previous word
			builder.WriteRune(previousRune)
			wordLength++
		}

		previousRune = r
	}

	// Write the last rune
	if previousRune != 0 {
		builder.WriteRune(previousRune)
	}

	return builder.String()
}

func extractSummary(doc string) string {
	firstParagraph := normalizeDoc(doc)

	// Trim to first paragraph if there are multiple
	if idx := strings.Index(firstParagraph, "\n\n"); idx != -1 {
		firstParagraph = firstParagraph[:idx]
	}

	var builder strings.Builder
	var lookback1 rune
	var lookback2 rune
	var lookback3 rune
	for _, r := range formatDocParagraph(firstParagraph) {
		// We terminate the sequence if we see a space preceded by a '.' which
		// does not have exactly one word character before it (to avoid
		// treating initials as the end of a sentence).
		isPeriod := r == ' ' && lookback1 == '.'
		isInitial := unicode.IsUpper(lookback2) && !unicode.IsLetter(lookback3) && !unicode.IsDigit(lookback3)
		if isPeriod && !isInitial {
			break
		}

		// Write the rune
		builder.WriteRune(r)

		// Update tracking variables
		lookback3 = lookback2
		lookback2 = lookback1
		lookback1 = r
	}

	// Make the summary end with a period if it is nonempty and doesn't already.
	if lookback1 != '.' && lookback1 != 0 {
		builder.WriteRune('.')
	}

	return builder.String()
}

var crlfRegex = regexp.MustCompile("\r\n")

func normalizeDoc(doc string) string {
	return strings.TrimSpace(crlfRegex.ReplaceAllString(doc, "\n"))
}

func formatDocParagraph(paragraph string) string {
	var mergedParagraph strings.Builder
	for i, line := range strings.Split(paragraph, "\n") {
		if i > 0 {
			mergedParagraph.WriteRune(' ')
		}

		mergedParagraph.WriteString(strings.TrimSpace(line))
	}

	return mergedParagraph.String()
}

func createDeclCopyWithoutComments(from *ast.GenDecl) *ast.GenDecl {
	var specs []ast.Spec
	for _, spec := range from.Specs {
		switch spec.(type) {
		case *ast.TypeSpec:
			typeSpec := spec.(*ast.TypeSpec)
			switch typeSpec.Type.(type) {
			case *ast.StructType:
				structType := *typeSpec.Type.(*ast.StructType)

				var copyFields []*ast.Field
				fields := structType.Fields

				for _, field := range fields.List {
					fieldCopy := copyFieldWithoutDoc(field)
					copyFields = append(copyFields, fieldCopy)
				}

				listCopy := copyFieldListWithFields(fields, copyFields)
				structTypeCopy := copyStructTypeWithFieldList(structType, listCopy)
				specs = append(specs, copySpecWithStructType(typeSpec, structTypeCopy))
			}
		}
	}

	return copyGenDeclWithSpecs(from, specs)
}

func copyGenDeclWithSpecs(decl *ast.GenDecl, specs []ast.Spec) *ast.GenDecl {
	return &ast.GenDecl{
		Doc:    decl.Doc,
		TokPos: decl.TokPos,
		Tok:    decl.Tok,
		Lparen: decl.Lparen,
		Specs:  specs,
		Rparen: decl.Rparen,
	}
}

func copySpecWithStructType(spec *ast.TypeSpec, typeSpec *ast.StructType) *ast.TypeSpec {
	return &ast.TypeSpec{
		Doc:        spec.Doc,
		Name:       spec.Name,
		TypeParams: spec.TypeParams,
		Assign:     spec.Assign,
		Type:       typeSpec,
		Comment:    spec.Comment,
	}
}

func copyStructTypeWithFieldList(structTyp ast.StructType, fieldList *ast.FieldList) *ast.StructType {
	return &ast.StructType{
		Struct:     structTyp.Struct,
		Fields:     fieldList,
		Incomplete: structTyp.Incomplete,
	}
}

func copyFieldListWithFields(list *ast.FieldList, fields []*ast.Field) *ast.FieldList {
	return &ast.FieldList{
		Opening: list.Opening,
		List:    fields,
		Closing: list.Closing,
	}
}

func copyFieldWithoutDoc(field *ast.Field) *ast.Field {
	return &ast.Field{
		Doc:     nil,
		Names:   field.Names,
		Type:    field.Type,
		Tag:     field.Tag,
		Comment: field.Comment,
	}
}
