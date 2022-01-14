package super

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Rid/hapesay/cmd/v2/internal/screen"
	hapesay "github.com/Rid/hapesay/v2"
	"github.com/Rid/hapesay/v2/decoration"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

func getNoSaidHape(hape *hapesay.Hape, opts ...hapesay.Option) (string, error) {
	opts = append(opts, hapesay.Thoughts1(' '))
	opts = append(opts, hapesay.Thoughts2(' '))

	hape, err := hape.Clone(opts...)
	if err != nil {
		return "", err
	}
	return hape.GetHape()
}

// RunSuperHape runs super hape mode animation on the your terminal
func RunSuperHape(phrase string, withBold bool, opts ...hapesay.Option) error {
	hape, err := hapesay.New(opts...)
	if err != nil {
		return err
	}
	balloon := hape.Balloon(phrase)
	blank := createBlankSpace(balloon)

	said, err := hape.GetHape()
	if err != nil {
		return err
	}

	notSaid, err := getNoSaidHape(hape, opts...)
	if err != nil {
		return err
	}

	saidHape := balloon + said
	saidHapeLines := strings.Count(saidHape, "\n") + 1

	// When it is higher than the height of the terminal
	h := screen.Height()
	if saidHapeLines > h {
		return errors.New("too height messages")
	}

	notSaidHape := blank + notSaid

	renderer := newRenderer(saidHape, notSaidHape)

	screen.SaveState()
	screen.HideCursor()
	screen.Clear()

	go renderer.createFrames(hape, withBold)

	renderer.render()

	screen.UnHideCursor()
	screen.RestoreState()

	return nil
}

func createBlankSpace(balloon string) string {
	var buf strings.Builder
	l := strings.Count(balloon, "\n")
	for i := 0; i < l; i++ {
		buf.WriteRune('\n')
	}
	return buf.String()
}

func maxLen(hape []string) int {
	max := 0
	for _, line := range hape {
		l := runewidth.StringWidth(line)
		if max < l {
			max = l
		}
	}
	return max
}

type hapeLine struct {
	raw      string
	clusters []rune
}

func (c *hapeLine) Len() int {
	return len(c.clusters)
}

func (c *hapeLine) Slice(i, j int) string {
	if c.Len() == 0 {
		return ""
	}
	return string(c.clusters[i:j])
}

func makeHapeLines(hape string) []*hapeLine {
	sep := strings.Split(hape, "\n")
	hapeLines := make([]*hapeLine, len(sep))
	for i, line := range sep {
		g := uniseg.NewGraphemes(line)
		clusters := make([]rune, 0)
		for g.Next() {
			clusters = append(clusters, g.Runes()...)
		}
		hapeLines[i] = &hapeLine{
			raw:      line,
			clusters: clusters,
		}
	}
	return hapeLines
}

type renderer struct {
	max         int
	middle      int
	screenWidth int
	heightDiff  int
	frames      chan string

	saidHape         string
	notSaidHapeLines []*hapeLine

	quit chan os.Signal
}

func newRenderer(saidHape, notSaidHape string) *renderer {
	notSaidHapeSep := strings.Split(notSaidHape, "\n")
	w, hapesWidth := screen.Width(), maxLen(notSaidHapeSep)
	max := w + hapesWidth

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return &renderer{
		max:              max,
		middle:           max / 2,
		screenWidth:      w,
		heightDiff:       screen.Height() - strings.Count(saidHape, "\n") - 1,
		frames:           make(chan string, max),
		saidHape:         saidHape,
		notSaidHapeLines: makeHapeLines(notSaidHape),
		quit:             quit,
	}
}

const (
	// Frequency the color changes
	magic = 2

	span    = 30 * time.Millisecond
	standup = 3 * time.Second
)

func (r *renderer) createFrames(hape *hapesay.Hape, withBold bool) {
	const times = standup / span
	w := r.newWriter(withBold)

	for x, i := 0, 1; i <= r.max; i++ {
		if i == r.middle {
			w.SetPosx(r.posX(i))
			for k := 0; k < int(times); k++ {
				base := x * 70
				// draw colored hape
				w.SetColorSeq(base)
				w.WriteString(r.saidHape)
				r.frames <- w.String()
				if k%magic == 0 {
					x++
				}
			}
		} else {
			base := x * 70
			w.SetPosx(r.posX(i))
			w.SetColorSeq(base)

			for _, line := range r.notSaidHapeLines {
				if i > r.screenWidth {
					// Left side animations
					n := i - r.screenWidth
					if n < line.Len() {
						w.WriteString(line.Slice(n, line.Len()))
					}
				} else if i <= line.Len() {
					// Right side animations
					w.WriteString(line.Slice(0, i-1))
				} else {
					w.WriteString(line.raw)
				}
				w.Write([]byte{'\n'})
			}
			r.frames <- w.String()
		}
		if i%magic == 0 {
			x++
		}
	}
	close(r.frames)
}

func (r *renderer) render() {
	initCh := make(chan struct{}, 1)
	initCh <- struct{}{}

	for view := range r.frames {
		select {
		case <-r.quit:
			screen.Clear()
			return
		case <-initCh:
		case <-time.After(span):
		}
		io.Copy(screen.Stdout, strings.NewReader(view))
	}
}

func (r *renderer) posX(i int) int {
	posx := r.screenWidth - i
	if posx < 1 {
		posx = 1
	}
	return posx
}

// Writer is wrapper which is both screen.MoveWriter and decoration.Writer.
type Writer struct {
	buf *strings.Builder
	mw  *screen.MoveWriter
	dw  *decoration.Writer
}

func (r *renderer) newWriter(withBold bool) *Writer {
	var buf strings.Builder
	mw := screen.NewMoveWriter(&buf, r.posX(0), r.heightDiff)
	options := []decoration.Option{
		decoration.WithAurora(0),
	}
	if withBold {
		options = append(options, decoration.WithBold())
	}
	dw := decoration.NewWriter(mw, options...)
	return &Writer{
		buf: &buf,
		mw:  mw,
		dw:  dw,
	}
}

// WriteString writes string. which is implemented io.StringWriter.
func (w *Writer) WriteString(s string) (int, error) { return w.dw.WriteString(s) }

// Write writes bytes. which is implemented io.Writer.
func (w *Writer) Write(p []byte) (int, error) { return w.dw.Write(p) }

// SetPosx sets posx
func (w *Writer) SetPosx(x int) { w.mw.SetPosx(x) }

// SetColorSeq sets color sequence.
func (w *Writer) SetColorSeq(colorSeq int) { w.dw.SetColorSeq(colorSeq) }

// Reset resets calls some Reset methods.
func (w *Writer) Reset() {
	w.buf.Reset()
	w.mw.Reset()
}

func (w *Writer) String() string {
	defer w.Reset()
	return w.buf.String()
}
