---
title: Books
---

{{ range .Data }}
{{ range .Data.book }}
## {{ .title }}

![Thumbnail]({{ .thumbnail }})

Quality: {{ .quality }}

Completion: {{ .completion }}

Reading: {{ .reading }}

Authors:
{{ range .authors }}
- {{ .main }} {{ if .alts }}({{ range .alts }}{{ . }}{{ end }}){{ end }}
{{ end }}

Links:
{{ range .links }}
- {{ .url }}
{{ end }}

{{ end }}
{{ end }}
