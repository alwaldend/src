---
title: Books
---

<table>
  <caption>
    Books
  </caption>
  <thead>
    <tr>
      <th scope="col">Title</th>
      <th scope="col">Info</th>
    </tr>
  </thead>
  <tbody>
    {{ range .Data -}}
    {{ range .Data.book -}}
    <tr>
      <td>
        {{ .title }}
        {{- range .alt_titles }} <br> ({{ . }}){{ end }}
        <br>
        <img src="{{ .thumbnail }}" alt="Thumbnail of {{ .title }}" height="300">
      </td>
      <td>
        {{ range .authors -}}
        Author: {{ .main }}{{ range .alts }} ({{ . }}){{ end }}
        <br>
        {{ end -}}
        Completion status: {{ .completion }}
        <br>
        Reading status: {{ .reading }}
        <br>
        Quality: {{ .quality }}
        {{ if .summary }}
        <br>
        Summary: {{ .summary }}
        {{ end -}}
        <br>
        {{ range .links -}}
        Link: <a href="{{ .url }}">{{ .title }}</a>
        <br>
        {{- end }}
      </td>
    </tr>
    {{ end -}}
    {{ end -}}
  </tbody>
</table>
