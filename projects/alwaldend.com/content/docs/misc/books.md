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
    <table class="td-initial table col text-left">
      <tbody>
        <tr>
          <td>Title</td>
          <td>{{ .title }}</td>
        </tr>
        {{ if .alt_titles -}}
        <tr>
          <td>ALternative titles</td>
          <td>
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
          <td>Authors</td>
          <td>
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
          <td>Completion status</td>
          <td>{{ .completion }}</td>
        </tr>
        {{ end -}}
        {{ if .reading -}}
        <tr>
          <td>Reading status</td>
          <td>{{ .reading }}</td>
        </tr>
        {{ end -}}
        {{ if .quality -}}
        <tr>
          <td>Quality</td>
          <td>{{ .quality }}</td>
        </tr>
        {{ end -}}
        {{ if .summary -}}
        <tr>
          <td>Summary</td>
          <td>{{ .summary }}</td>
        </tr>
        {{ end -}}
        {{ if .links -}}
        <tr>
          <td>Links</td>
          <td>
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
