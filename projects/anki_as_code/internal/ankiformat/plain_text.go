package ankiformat

import (
	"strings"

	"golang.org/x/net/html"
)

// isMediaTag identifies HTML elements whose filenames belong in the field cache.
func isMediaTag(name string) bool {
	switch name {
	case "img", "audio", "video", "object", "source":
		return true
	default:
		return false
	}
}

// PlainText extracts field-cache text from HTML while retaining media filenames.
func PlainText(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	var out strings.Builder
	hidden := 0
	for {
		t := z.Next()
		if t == html.ErrorToken {
			break
		}
		token := z.Token()
		if (t == html.StartTagToken || t == html.SelfClosingTagToken) && isMediaTag(token.Data) {
			for _, a := range token.Attr {
				if a.Key == "src" || a.Key == "data" {
					out.WriteString(" ")
					out.WriteString(a.Val)
					out.WriteString(" ")
					break
				}
			}
		}
		if t == html.StartTagToken && (token.Data == "script" || token.Data == "style") {
			hidden++
		}
		if t == html.TextToken && hidden == 0 {
			out.WriteString(token.Data)
		}
		if t == html.EndTagToken && (token.Data == "script" || token.Data == "style") && hidden > 0 {
			hidden--
		}
	}
	return strings.ReplaceAll(out.String(), "\u00a0", " ")
}
