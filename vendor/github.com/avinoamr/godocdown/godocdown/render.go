package main

import (
	"fmt"
	"go/doc"
	"io"
)

func renderConstantSectionTo(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), filterText(entry.Doc))
	}
}

func renderVariableSectionTo(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), filterText(entry.Doc))
	}
}

func renderFunctionSectionTo(writer io.Writer, list []*doc.Func, inTypeSection bool, exs []*doc.Example) {

	header := RenderStyle.FunctionHeader
	if inTypeSection {
		header = RenderStyle.TypeFunctionHeader
	}

	for _, entry := range list {
		receiver := " "
		if entry.Recv != "" {
			receiver = fmt.Sprintf("(%s) ", entry.Recv)
		}
		fmt.Fprintf(writer, "%s <a name='%s'></a> func %s[%s]()\n\n%s\n%s\n",
			header,
			entry.Name,
			receiver,
			entry.Name,
			indentCode(sourceOfNode(entry.Decl)),
			filterText(entry.Doc)) // use the doc as-is in markdown

		for _, ex := range filterExamples(exs, entry.Name) {
			renderExample(writer, ex)
		}
	}
}

func renderExample(w io.Writer, ex *doc.Example) {
	code := sourceOfNode(ex.Code)
	code = indentCode(code)

	_, sub := exampleNames(ex.Name)
	fmt.Fprintf(w, "<a name='Example%s'></a><details><summary>Example%s</summary><p>\n\n%s\n%s\n\nOutput:\n```\n%s```\n</p></details>\n\n",
		ex.Name,
		sub,
		filterText(ex.Doc),
		code,
		ex.Output)
}

func renderTypeSectionTo(writer io.Writer, list []*doc.Type, exs []*doc.Example) {
	header := RenderStyle.TypeHeader

	for _, entry := range list {
		fmt.Fprintf(writer, "%s <a name='%s'></a>type [%s]()\n\n%s\n\n%s\n",
			header,
			entry.Name,
			entry.Name,
			indentCode(sourceOfNode(entry.Decl)),
			filterText(entry.Doc))

		for _, ex := range filterExamples(exs, entry.Name) {
			renderExample(writer, ex)
		}

		renderConstantSectionTo(writer, entry.Consts)
		renderVariableSectionTo(writer, entry.Vars)
		renderFunctionSectionTo(writer, entry.Funcs, true, exs)
		renderFunctionSectionTo(writer, entry.Methods, true, nil)
	}
}

func renderHeaderTo(writer io.Writer, document *_document) {
	fmt.Fprintf(writer, "# %s\n\n", document.Name)

	if !document.IsCommand {
		// Import
		if RenderStyle.IncludeImport {
			if document.ImportPath != "" {
				code := fmt.Sprintf(`import "%s"`, document.ImportPath)
				code = indentCode(code)
				fmt.Fprintf(writer, "%s\n\n", code)
			}
		}
	}
}

func renderSynopsisTo(writer io.Writer, document *_document) {
	fmt.Fprintf(writer, "%s\n", headifySynopsis(filterText(document.pkg.Doc)))
}

func renderUsageTo(writer io.Writer, document *_document) {

	exs := document.Examples

	// Usage
	fmt.Fprintf(writer, "%s\n", RenderStyle.UsageHeader)

	// render index
	renderIndex(writer, document, exs)

	// Constant Section
	renderConstantSectionTo(writer, document.pkg.Consts)

	// Variable Section
	renderVariableSectionTo(writer, document.pkg.Vars)

	// Function Section
	renderFunctionSectionTo(writer, document.pkg.Funcs, false, exs)

	// Type Section
	renderTypeSectionTo(writer, document.pkg.Types, exs)
}

func renderSignatureTo(writer io.Writer) {
	if RenderStyle.IncludeSignature {
		fmt.Fprintf(writer, "\n\n--\n**godocdown** http://github.com/avinoamr/godocdown\n")
	}
}

func renderFunctionIndexTo(w io.Writer, list []*doc.Func, inType bool) {
	prefix := ""
	if inType {
		prefix = "    "
	}

	for _, e := range list {
		decl := sourceOfNode(e.Decl)
		fmt.Fprintf(w, "%s - [%s](#%s)\n", prefix, decl, e.Name)
	}
}

func renderTypeIndexTo(w io.Writer, list []*doc.Type) {
	for _, e := range list {
		fmt.Fprintf(w, " - [type %s](#%s)\n", e.Name, e.Name)
		renderFunctionIndexTo(w, e.Funcs, true)
	}
}

func renderExampleIndexTo(w io.Writer, list []*doc.Example) {
	if len(list) == 0 {
		return
	}

	fmt.Fprintf(w, "\n#### Examples\n\n")
	for _, e := range list {
		name, sub := exampleNames(e.Name)
		fmt.Fprintf(w, " - [%s%s](#Example%s)\n", name, sub, e.Name)
	}
}

func renderIndex(w io.Writer, d *_document, exs []*doc.Example) {
	renderFunctionIndexTo(w, d.pkg.Funcs, false)
	renderTypeIndexTo(w, d.pkg.Types)
	renderExampleIndexTo(w, exs)
	fmt.Fprintf(w, "\n")
}
