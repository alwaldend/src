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
        {{- if .alt_titles }}{{ range .alt_titles }} ({{ . }}){{ end }}<br>{{ end }}
        {{- if .authors }}{{ range .authors }} [{{ .main }}{{ range .alts }}, {{ . }}{{ end }}]{{ end }}<br>{{ end }}
        <img src="{{ .thumbnail }}" alt="Thumbnail of {{ .title }}" height="300"></img>
      </th>
      <td>{{ .quality }}</td>
      <td>{{ .completion }}</td>
      <td>{{ .reading }}</td>
      <td>
        <ul>
        {{ range .links -}}
        <li><a href="{{ .url }}">{{ .title }}</a></li>
        {{ end -}}
        </ul>
      </td>
    </tr>
    {{ end -}}
    {{ end -}}
  </tbody>
</table>
