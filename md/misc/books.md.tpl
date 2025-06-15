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
      <th scope="col">Links</th>
    </tr>
  </thead>
  <tbody>
    {{ range .Data -}}
    {{ range .Data.book -}}
    <tr>
      <th scope="row">
        {{ .title }}
        {{- range .alt_titles }} <br> ({{ . }}){{ end }}
        {{- range .authors }} <br> [{{ .main }}{{ range .alts }}, {{ . }}{{ end }}]{{ end }}
        <br>
        <img src="{{ .thumbnail }}" alt="Thumbnail of {{ .title }}" height="300"></img>
      </th>
      <td>
        Completion status: {{ .completion }}
        <br>
        Reading status: {{ .reading }}
        <br>
        Quality: {{ .quality }}
        <br>
        Summary: {{ .summary }}
      </td>
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
