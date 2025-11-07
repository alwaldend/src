---
title: Books
description: Books
tags:
  - books
---

## Books

{{< docs_misc_books.inline >}}

{{ $authors := (index .Site.Data.projects "alwaldend.com").authors.author -}}

<div class="container">
  <div class="row table-responsive">
    <table class="td-initial table text-start">
      <tbody>
        {{ range ((index .Site.Data.projects "alwaldend.com").books.book | collections.Reverse) -}}
        {{ $rowspan := 0 }}
        {{ range (collections.Slice .title .alt_titles .authors .completion .reading .quality .summary .links) }}
        {{ if . }}{{ $rowspan = math.Add $rowspan 1 }}{{ end }}
        {{ end }}
        <tr>
          <td colspan="1" class="text-center" rowspan="{{ $rowspan }}">
            <img
              src="{{ .thumbnail }}"
              class="img-fluid mx-auto"
              alt="Thumbnail of {{ .title }}">
          </td>
          <td colspan="1">Title</td>
          <td colspan="4">{{ printf "##### %s" .title | $.Page.RenderString }}</td>
        </tr>
        {{ if .alt_titles -}}
        <tr>
          <td colspan="1">ALternative titles</td>
          <td colspan="4">
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
          <td colspan="4">
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
          <td colspan="4">{{ .completion }}</td>
        </tr>
        {{ end -}}
        {{ if .reading -}}
        <tr>
          <td colspan="1">Reading status</td>
          <td colspan="4">{{ .reading }}</td>
        </tr>
        {{ end -}}
        {{ if .quality -}}
        <tr>
          <td colspan="1">Quality</td>
          <td colspan="4">{{ .quality }}</td>
        </tr>
        {{ end -}}
        {{ if .summary -}}
        <tr>
          <td colspan="1">Summary</td>
          <td colspan="4">{{ .summary }}</td>
        </tr>
        {{ end -}}
        {{ if .links -}}
        <tr>
          <td colspan="1">Links</td>
          <td colspan="4">
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
        {{ end -}}
      </tbody>
    </table>
  </div>
</div>


{{< /docs_misc_books.inline >}}
