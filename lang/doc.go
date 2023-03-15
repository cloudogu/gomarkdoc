package lang

import (
	"go/doc"
	"go/doc/comment"
)

// Doc provides access to the documentation comment contents for a package or
// symbol in a structured form.
type Doc struct {
	cfg            *Config
	blocks         []*Block
	currentPackage string
	types          []*doc.Type
}

// NewDoc initializes a Doc struct from the provided raw documentation text and
// with headers rendered by default at the heading level provided. Documentation
// is separated into block level elements using the standard rules from golang's
// documentation conventions.
func NewDoc(cfg *Config, text string) *Doc {
	return NewDocWithDocLinkParser(cfg, text, currentPackage, knownTypes)
}

// NewDocWithDocLinkParser initializes a Doc struct with additional information for modifying [comment.Parser].
// With the package and types the parser can parse [comment.DocLink].
func NewDocWithDocLinkParser(cfg *Config, text string, currentPackage string, types []*doc.Type) *Doc {
	// Replace CRLF with LF
	rawText := normalizeDoc(text)

	doc := Doc{cfg, nil, currentPackage, types}
	var p comment.Parser
	p.LookupPackage = doc.lookUpPackage
	p.LookupSym = doc.lookUpSymbol

	parsed := p.Parse(rawText)

	blocks := ParseBlocks(cfg, parsed.Content, false)
	doc.blocks = blocks
	return &doc
}

// lookUpPackage reports whether a package is the current package.
// Returns true if the package is currentPackage.
func (d *Doc) lookUpPackage(name string) (string, bool) {
	if d.currentPackage == "" {
		return "", false
	}

	if name == d.currentPackage {
		return d.currentPackage, true
	}

	return "", false
}

// lookUpSymbol reports whether a symbol name or method name exists in the current package.
// Returns true if the symbol is in types (a type in the actual package).
func (d *Doc) lookUpSymbol(_, name string) bool {
	if d.types == nil {
		return false
	}

	for _, e := range d.types {
		if e.Name == name {
			return true
		}
	}
	return false
}

// Level provides the default level that headers within the documentation should
// be rendered
func (d *Doc) Level() int {
	return d.cfg.Level
}

// Blocks holds the list of block elements that makes up the documentation
// contents.
func (d *Doc) Blocks() []*Block {
	return d.blocks
}
