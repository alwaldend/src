---
title: Books
---

{{< md_misc_books.inline >}}

<table>
  <caption>
    Books
  </caption>
  <thead>
    <tr>
      <th>Thumbnail</th>
      <th colspan="2">Info</th>
    </tr>
  </thead>
  <tbody>
    {{ $authors := .Site.Data.hugo.sites.docs.authors.author -}}
    {{ range .Site.Data.hugo.sites.docs.books.book -}}
    <tr>
      <td>
        <img
          src="{{ .thumbnail }}"
          alt="Thumbnail of {{ .title }}"
          width="300"
        >
      </td>
      <td colspan="2">
        <table>
          <tbody>
            <tr>
              <td>Title</td>
              <td>{{ .title }}</td>
            </tr>
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
                  <li>{{ .name }}{{ range .alts }} ({{ . }}){{ end }}</li>
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
      </td>
    </tr>
    {{ end -}}
  </tbody>
</table>

{{< /md_misc_books.inline >}}
