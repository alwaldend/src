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
      <th scope="row">
        {{ .title }}
        {{- range .alt_titles }} <br> ({{ . }}){{ end }}
        <br>
        <img src="{{ .thumbnail }}" alt="Thumbnail of {{ .title }}" height="300"></img>
      </th>
      <td>
        Authors:
        <ul>
          {{ range .authors -}}
          <li>
            {{ .main }}{{ range .alts }} ({{ . }}){{ end }}
          </li>
          {{ end -}}
        </ul>
        <br>
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
        Links:
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
