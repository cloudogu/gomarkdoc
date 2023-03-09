package lang

import (
	"fmt"
	"go/ast"
	"go/doc"
	"strings"
)

// Field holds documentation information for a single field declaration within a
// type.
type Field struct {
	cfg      *Config
	doc      *ast.Field
	examples []*doc.Example
}

// NewField creates a new Field from the corresponding documentation construct
// from the standard library, the related token.FileSet for the field and
// the list of examples for the field.
func NewField(cfg *Config, doc *ast.Field, examples []*doc.Example) *Field {
	return &Field{cfg, doc, examples}
}

// Level provides the default level at which headers for the field should be
// rendered in the final documentation.
func (f *Field) Level() int {
	return f.cfg.Level
}

// Name provides the name of the field.
func (f *Field) Name() string {
	return f.doc.Names[0].Name
}

// Title provides the formatted name of the field. It is primarily designed for
// generating headers.
func (f *Field) Title() string {
	return fmt.Sprintf("Field %s", f.Name())
}

// Summary provides the one-sentence summary of the field's documentation
// comment
func (f *Field) Summary() string {
	return extractSummary(f.doc.Doc.Text())
}

// Doc provides the structured contents of the documentation comment for the
// field.
func (f *Field) Doc() *Doc {
	return NewDoc(f.cfg.Inc(1), f.doc.Doc.Text())
}

// Decl provides the raw text representation of the code for declaring the const
// or var.
func (f *Field) Decl() (string, error) {
	return printNode(f.doc, f.cfg.FileSet)
}

// Examples provides the list of examples from the list given on initialization
// that pertain to the field.
func (f *Field) Examples() (examples []*Example) {
	fullName := f.Name()
	underscorePrefix := fmt.Sprintf("%s_", fullName)

	for _, example := range f.examples {
		var name string
		switch {
		case example.Name == fullName:
			name = ""
		case strings.HasPrefix(example.Name, underscorePrefix):
			name = example.Name[len(underscorePrefix):]
		default:
			// TODO: better filtering
			continue
		}

		examples = append(examples, NewExample(f.cfg.Inc(1), name, example))
	}

	return
}
