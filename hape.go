package hapesay

import (
	"fmt"
	"math/rand"
	"strings"
)

// Hape struct!!
type Hape struct {
	eyes            string
	tongue          string
	typ             *HapeFile
	thoughts        rune
	thinking        bool
	ballonWidth     int
	disableWordWrap bool
	balloonOffset   int

	buf strings.Builder
}

// New returns pointer of Hape struct that made by options
func New(options ...Option) (*Hape, error) {
	hape := &Hape{
		eyes:     "oo",
		tongue:   "  ",
		thoughts: '/',
		typ: &HapeFile{
			Name:         "default",
			BasePath:     "hapes",
			LocationType: InBinary,
		},
		ballonWidth: 40,
	}
	for _, o := range options {
		if err := o(hape); err != nil {
			return nil, err
		}
	}
	return hape, nil
}

// Say returns string that said by hape
func (hape *Hape) Say(phrase string) (string, error) {
	mow, err := hape.GetHape()
	if err != nil {
		return "", err
	}
	return hape.Balloon(phrase) + mow, nil
}

// Clone returns a copy of hape.
//
// If any options are specified, they will be reflected.
func (hape *Hape) Clone(options ...Option) (*Hape, error) {
	ret := new(Hape)
	*ret = *hape
	ret.buf.Reset()
	for _, o := range options {
		if err := o(ret); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// Option defined for Options
type Option func(*Hape) error

// Eyes specifies eyes
// The specified string will always be adjusted to be equal to two characters.
func Eyes(s string) Option {
	return func(c *Hape) error {
		c.eyes = adjustTo2Chars(s)
		return nil
	}
}

// Tongue specifies tongue
// The specified string will always be adjusted to be less than or equal to two characters.
func Tongue(s string) Option {
	return func(c *Hape) error {
		c.tongue = adjustTo2Chars(s)
		return nil
	}
}

func adjustTo2Chars(s string) string {
	if len(s) >= 2 {
		return s[:2]
	}
	if len(s) == 1 {
		return s + " "
	}
	return "  "
}

func containHapes(target string) (*HapeFile, error) {
	hapePaths, err := Hapes()
	if err != nil {
		return nil, err
	}
	for _, hapePath := range hapePaths {
		hapefile, ok := hapePath.Lookup(target)
		if ok {
			return hapefile, nil
		}
	}
	return nil, nil
}

// NotFound is indicated not found the hapefile.
type NotFound struct {
	Hapefile string
}

var _ error = (*NotFound)(nil)

func (n *NotFound) Error() string {
	return fmt.Sprintf("not found %q hapefile", n.Hapefile)
}

// Type specify name of the hapefile
func Type(s string) Option {
	if s == "" {
		s = "default"
	}
	return func(c *Hape) error {
		hapefile, err := containHapes(s)
		if err != nil {
			return err
		}
		if hapefile != nil {
			c.typ = hapefile
			return nil
		}
		return &NotFound{Hapefile: s}
	}
}

// Thinking enables thinking mode
func Thinking() Option {
	return func(c *Hape) error {
		c.thinking = true
		return nil
	}
}

// Thoughts Thoughts allows you to specify
// the rune that will be drawn between
// the speech bubbles and the hape
func Thoughts(thoughts rune) Option {
	return func(c *Hape) error {
		c.thoughts = thoughts
		return nil
	}
}

// Random specifies something .hape from hapes directory
func Random() Option {
	pick, err := pickHape()
	return func(c *Hape) error {
		if err != nil {
			return err
		}
		c.typ = pick
		return nil
	}
}

func pickHape() (*HapeFile, error) {
	hapePaths, err := Hapes()
	if err != nil {
		return nil, err
	}
	hapePath := hapePaths[rand.Intn(len(hapePaths))]

	n := len(hapePath.HapeFiles)
	hapefile := hapePath.HapeFiles[rand.Intn(n)]
	return &HapeFile{
		Name:         hapefile,
		BasePath:     hapePath.Name,
		LocationType: hapePath.LocationType,
	}, nil
}

// BallonWidth specifies ballon size
func BallonWidth(size uint) Option {
	return func(c *Hape) error {
		c.ballonWidth = int(size)
		return nil
	}
}

// DisableWordWrap disables word wrap.
// Ignoring width of the ballon.
func DisableWordWrap() Option {
	return func(c *Hape) error {
		c.disableWordWrap = true
		return nil
	}
}
