---
title: Books
description: books
tags:
  - books
---

## Books


{{< docs_misc_books.inline >}}

{{ $authors := (index .Site.Data.projects "alwaldend.com").authors.author -}}
{{ range ((index .Site.Data.projects "alwaldend.com").books.book | collections.Reverse) -}}

<div class="container">
  <div class="row table-responsive">
    <img
      src="{{ .thumbnail }}"
      class="col-3"
      alt="Thumbnail of {{ .title }}">
    <table class="td-initial table col text-start">
      <tbody>
        <tr>
          <td colspan="1">Title</td>
          <td colspan="2">{{ .title }}</td>
        </tr>
        {{ if .alt_titles -}}
        <tr>
          <td colspan="1">ALternative titles</td>
          <td colspan="2">
            <ul>
              {{ range .alt_titles -}}
              <li>
                {{ . }}
              </li>
              {{- end }}
            </ul>
          </td>
        </tr>
        {{ end -}}
        {{ if .authors -}}
        <tr>
          <td colspan="1">Authors</td>
          <td colspan="2">
            <ul>
              {{ range .authors -}}
              {{ with index $authors . -}}
              <li>{{ .name }}{{ range .alt_names }} ({{ . }}){{ end }}</li>
              {{- end }}
              {{- end }}
            </ul>
          </td>
        </tr>
        {{ end -}}
        {{ if .completion -}}
        <tr>
          <td colspan="1">Completion status</td>
          <td colspan="2">{{ .completion }}</td>
        </tr>
        {{ end -}}
        {{ if .reading -}}
        <tr>
          <td colspan="1">Reading status</td>
          <td colspan="2">{{ .reading }}</td>
        </tr>
        {{ end -}}
        {{ if .quality -}}
        <tr>
          <td colspan="1">Quality</td>
          <td colspan="2">{{ .quality }}</td>
        </tr>
        {{ end -}}
        {{ if .summary -}}
        <tr>
          <td colspan="1">Summary</td>
          <td colspan="2">{{ .summary }}</td>
        </tr>
        {{ end -}}
        {{ if .links -}}
        <tr>
          <td colspan="1">Links</td>
          <td colspan="2">
            <ul>
              {{ range .links -}}
              <li>
                <a href="{{ .url }}">
                  {{ .title }}
                </a>
              </li>
              {{- end }}
            </ul>
          </td>
        </tr>
        {{- end }}
      </tbody>
    </table>
  </div>
</div>

{{ end -}}

{{< /docs_misc_books.inline >}}
