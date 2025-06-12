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
      <th scope="col">Authors</th>
      <th scope="col">Quality</th>
      <th scope="col">Completion</th>
      <th scope="col">Reading</th>
      <th scope="col">Links</th>
    </tr>
  </thead>
  <tbody>
    {{ range .Data -}}
    {{ range .Data.book -}}
    <tr>
      <th scope="row">
        {{ .title }}
        <br>
        <img src="{{ .thumbnail }}" height=300>
      </th>
      <td>
        <ul>
        {{ range .authors -}}
        <li>{{ .main }} {{ if .alts }}({{ range .alts }}{{ . }}{{ end }}){{ end }}</li>
        {{ end -}}
        </ul>
      </td>
      <td>{{ .quality }}</td>
      <td>{{ .completion }}</td>
      <td>{{ .reading }}</td>
      <td>
        <ul>
        {{ range .links -}}
        <li><a href="{{ .url }}">{{ if .title }}{{ . }}{{ else }}{{ .url }}{{ end }}</a></li>
        {{ end -}}
        </ul>
      </td>
    </tr>
    {{ end -}}
    {{ end -}}
  </tbody>
</table>
