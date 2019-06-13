package main

const tplTxt = `
{{with $examples := .Examples}}
{{define "Code"}}
{{indentNode .Decl}}
{{filterText .Doc}}
{{end}}

{{define "Func"}}
## <a name="{{.Name}}"></a> func {{if .Recv}}({{.Recv}}){{end}}[{{.Name}}](#{{.Name}})
{{template "Code" .}}

{{range $examples}}{{template "Example" .}}{{end}}
{{end}}

{{define "Example"}}
<a name="Example{{.Name}}"></a><details><summary>Example{{exampleSubName .Name}}</summary><p>
{{filterText .Doc}}
{{indentNode .Code}}

Output:
{{indentCode .Output}}
{{end}}

{{define "Type"}}
## <a name="{{.Name}}"></a> type [{{.Name}}](#{{.Name}})
{{template "Code" .}}
{{range .Consts}}{{template "Code" .}}{{end}}
{{range .Vars}}{{template "Code" .}}{{end}}
{{range .Funcs}}{{template "Func" .}}{{end}}
{{range .Methods}}{{template "Func" .}}{{end}}
{{end}}

# {{.Name}} {{.Badge}}

{{indentCode .Import}}

{{.Synopsis}}
{{range .Consts}}{{template "Code" .}}{{end}}
{{range .Vars}}{{template "Code" .}}{{end}}
{{range .Funcs}}{{template "Func" .}}{{end}}
{{range .Types}}{{template "Type" .}}{{end}}
`
