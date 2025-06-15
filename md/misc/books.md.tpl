---
title: Books
---

<table>
  <caption>
    Books
  </caption>
  <thead>
    <tr>
      <th>Title</th>
      <th>Info</th>
    </tr>
  </thead>
  <tbody>
    {{ range .Data -}}
    {{ range .Data.book -}}
    <tr>
      <td>
        {{ .title }}
        <br>
        <img
          src="{{ .thumbnail }}"
          alt="Thumbnail of {{ .title }}"
          width="400"
        >
      </td>
      <td>
        <table>
          <tbody>
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
                  <li>{{ .main }}{{ range .alts }} ({{ . }}){{ end }}</li>
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
    {{ end -}}
  </tbody>
</table>
