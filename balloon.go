package hapesay

import (
	"fmt"
	"strings"

	wordwrap "github.com/Code-Hex/go-wordwrap"
	runewidth "github.com/mattn/go-runewidth"
)

type border struct {
	first  [2]rune
	middle [2]rune
	last   [2]rune
	only   [2]rune
}

func (hape *Hape) borderType() border {
	if hape.thinking {
		return border{
			first:  [2]rune{'(', ')'},
			middle: [2]rune{'(', ')'},
			last:   [2]rune{'(', ')'},
			only:   [2]rune{'(', ')'},
		}
	}

	return border{
		first:  [2]rune{'/', '\\'},
		middle: [2]rune{'|', '|'},
		last:   [2]rune{'\\', '/'},
		only:   [2]rune{'<', '>'},
	}
}

type line struct {
	text      string
	runeWidth int
}

type lines []*line

func (hape *Hape) maxLineWidth(lines []*line) int {
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
		if !hape.disableWordWrap && maxWidth > hape.ballonWidth {
			return hape.ballonWidth
		}
	}
	return maxWidth
}

func (hape *Hape) getLines(phrase string) []*line {
	text := hape.canonicalizePhrase(phrase)
	lineTexts := strings.Split(text, "\n")
	lines := make([]*line, 0, len(lineTexts))
	for _, lineText := range lineTexts {
		lines = append(lines, &line{
			text:      lineText,
			runeWidth: runewidth.StringWidth(lineText),
		})
	}
	return lines
}

func (hape *Hape) canonicalizePhrase(phrase string) string {
	// Replace tab to 8 spaces
	phrase = strings.Replace(phrase, "\t", "       ", -1)

	if hape.disableWordWrap {
		return phrase
	}
	width := hape.ballonWidth
	return wordwrap.WrapString(phrase, uint(width))
}

// Balloon to get the balloon and the string entered in the balloon.
func (hape *Hape) Balloon(phrase string) string {
	defer hape.buf.Reset()

	lines := hape.getLines(phrase)
	maxWidth := hape.maxLineWidth(lines)

	hape.writeBallon(lines, maxWidth)

	return hape.buf.String()
}

func (hape *Hape) writeBallon(lines []*line, maxWidth int) {
	top := make([]byte, 0)
	bottom := make([]byte, 0)

	for i := 0; i < hape.balloonOffset; i++ {
		top = append(top, ' ')
		bottom = append(bottom, ' ')
	}

	for i := 0; i < maxWidth+2; i++ {
		top = append(top, '_')
		bottom = append(bottom, '-')
	}

	borderType := hape.borderType()

	hape.buf.Write(top)
	hape.buf.Write([]byte{' ', '\n'})
	defer func() {
		hape.buf.Write(bottom)
		hape.buf.Write([]byte{' ', '\n'})
	}()

	l := len(lines)
	if l == 1 {
		border := borderType.only
		for i := 0; i < (hape.balloonOffset - 1); i++ {
			hape.buf.WriteRune(' ')
		}
		hape.buf.WriteRune(border[0])
		hape.buf.WriteRune(' ')
		hape.buf.WriteString(lines[0].text)
		hape.buf.WriteRune(' ')
		hape.buf.WriteRune(border[1])
		hape.buf.WriteRune('\n')
		return
	}

	var border [2]rune
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		for i := 0; i < (hape.balloonOffset - 1); i++ {
			hape.buf.WriteRune(' ')
		}
		hape.buf.WriteRune(border[0])
		hape.buf.WriteRune(' ')
		hape.padding(lines[i], maxWidth)
		hape.buf.WriteRune(' ')
		hape.buf.WriteRune(border[1])
		hape.buf.WriteRune('\n')
	}
}

func (hape *Hape) flush(text, top, bottom fmt.Stringer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text.String(),
		bottom.String(),
	)
}

func (hape *Hape) padding(line *line, maxWidth int) {
	if maxWidth <= line.runeWidth {
		hape.buf.WriteString(line.text)
		return
	}

	hape.buf.WriteString(line.text)
	l := maxWidth - line.runeWidth
	for i := 0; i < l; i++ {
		hape.buf.WriteRune(' ')
	}
}
