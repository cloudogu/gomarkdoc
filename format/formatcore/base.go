package formatcore

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
	"mvdan.cc/xurls/v2"
)

// Bold converts the provided text to bold
func Bold(text string) string {
	if text == "" {
		return ""
	}

	return fmt.Sprintf("**%s**", Escape(text))
}

// CodeBlock wraps the provided code as a code block. Language syntax
// highlighting is not supported.
func CodeBlock(code string) string {
	var builder strings.Builder

	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if i != 0 {
			builder.WriteRune('\n')
		}

		builder.WriteRune('\t')
		builder.WriteString(line)
	}

	return builder.String()
}

// GFMCodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided), using
// the triple backtick format from GitHub Flavored Markdown.
func GFMCodeBlock(language, code string) string {
	return fmt.Sprintf("```%s\n%s\n```", language, strings.TrimSpace(code))
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func Header(level int, text string) (string, error) {
	if level < 1 {
		return "", errors.New("format: header level cannot be less than 1")
	}

	switch level {
	case 1:
		return fmt.Sprintf("# %s", text), nil
	case 2:
		return fmt.Sprintf("## %s", text), nil
	case 3:
		return fmt.Sprintf("### %s", text), nil
	case 4:
		return fmt.Sprintf("#### %s", text), nil
	case 5:
		return fmt.Sprintf("##### %s", text), nil
	default:
		// Only go up to 6 levels. Anything higher is also level 6
		return fmt.Sprintf("###### %s", text), nil
	}
}

// Link generates a link with the given text and href values.
func Link(text, href string) string {
	if text == "" {
		return ""
	}

	if href == "" {
		return text
	}

	return fmt.Sprintf("[%s](<%s>)", text, href)
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func ListEntry(depth int, text string) string {
	// TODO: this is a weird special case
	if text == "" {
		return ""
	}

	prefix := strings.Repeat("  ", depth)
	return fmt.Sprintf("%s- %s", prefix, text)
}

// GFMAccordion generates a collapsible content. The accordion's visible title
// while collapsed is the provided title and the expanded content is the body.
func GFMAccordion(title, body string) string {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>%s</p>\n</details>", title, Escape(body))
}

// GFMAccordionHeader generates the header visible when an accordion is
// collapsed.
//
// The GFMAccordionHeader is expected to be used in conjunction with
// GFMAccordionTerminator() when the demands of the body's rendering requires
// it to be generated independently. The result looks conceptually like the
// following:
//
//	accordion := GFMAccordionHeader("Accordion Title") + "Accordion Body" + GFMAccordionTerminator()
func GFMAccordionHeader(title string) string {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>", title)
}

// GFMAccordionTerminator generates the code necessary to terminate an
// accordion after the body. It is expected to be used in conjunction with
// GFMAccordionHeader(). See GFMAccordionHeader for a full description.
func GFMAccordionTerminator() string {
	return "</p>\n</details>"
}

// Paragraph formats a paragraph with the provided text as the contents
func Paragraph(text string) string {
	return text
}

var (
	specialCharacterRegex = regexp.MustCompile("([\\\\`*_{}\\[\\]()<>#+\\-!~])")
	urlRegex              = xurls.Strict() // Require a scheme in URLs
)

// Escape escapes the special characters in the provided text, but leaves URLs
// found intact. Note that the URLs included must begin with a scheme to skip
// the escaping.
func Escape(text string) string {
	b := []byte(text)

	var (
		cursor  int
		builder strings.Builder
	)

	for _, urlLoc := range urlRegex.FindAllIndex(b, -1) {
		// Walk through each found URL, escaping the text before the URL and
		// leaving the text in the URL unchanged.
		if urlLoc[0] > cursor {
			// Escape the previous section if its length is nonzero
			builder.Write(escapeRaw(b[cursor:urlLoc[0]]))
		}

		// Add the unescaped URL to the end of it
		builder.Write(b[urlLoc[0]:urlLoc[1]])

		// Move the cursor forward for the next iteration
		cursor = urlLoc[1]
	}

	// Escape the end of the string after the last URL if there's anything left
	if len(b) > cursor {
		builder.Write(escapeRaw(b[cursor:]))
	}

	return builder.String()
}

func escapeRaw(segment []byte) []byte {
	return specialCharacterRegex.ReplaceAll(segment, []byte("\\$1"))
}

// PlainText converts a markdown string to the plain text that appears in the
// rendered output.
func PlainText(text string) string {
	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	node := md.Parse([]byte(text))

	var builder strings.Builder
	plainTextInner(node, &builder)

	return builder.String()
}

func plainTextInner(node *blackfriday.Node, builder *strings.Builder) {
	// Only text nodes produce output
	if node.Type == blackfriday.Text {
		builder.Write(node.Literal)
	}

	// Run the children first
	if node.FirstChild != nil {
		plainTextInner(node.FirstChild, builder)
	}

	// Then run any other siblings
	if node.Next != nil {
		// Add extra space if necessary between nodes
		if node.Type == blackfriday.Paragraph ||
			node.Type == blackfriday.CodeBlock ||
			node.Type == blackfriday.Heading {
			builder.WriteRune(' ')
		}

		plainTextInner(node.Next, builder)
	}
}
