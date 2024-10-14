package md

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(content []byte) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(content, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
