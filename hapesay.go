package hapesay

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Say to return hapesay string.
func Say(phrase string, options ...Option) (string, error) {
	hape, err := New(options...)
	if err != nil {
		return "", err
	}
	return hape.Say(phrase)
}

// LocationType indicates the type of COWPATH.
type LocationType int

const (
	// InBinary indicates the COWPATH in binary.
	InBinary LocationType = iota

	// InDirectory indicates the COWPATH in your directory.
	InDirectory
)

// HapePath is information of the COWPATH.
type HapePath struct {
	// Name is name of the COWPATH.
	// If you specified `COWPATH=/foo/bar`, Name is `/foo/bar`.
	Name string
	// HapeFiles are name of the hapefile which are trimmed ".hape" suffix.
	HapeFiles []string
	// LocationType is the type of COWPATH
	LocationType LocationType
}

// Lookup will look for the target hapefile in the specified path.
// If it exists, it returns the hapefile information and true value.
func (c *HapePath) Lookup(target string) (*HapeFile, bool) {
	for _, hapefile := range c.HapeFiles {
		if hapefile == target {
			return &HapeFile{
				Name:         hapefile,
				BasePath:     c.Name,
				LocationType: c.LocationType,
			}, true
		}
	}
	return nil, false
}

// HapeFile is information of the hapefile.
type HapeFile struct {
	// Name is name of the hapefile.
	Name string
	// BasePath is the path which the hapepath is in.
	BasePath string
	// LocationType is the type of COWPATH
	LocationType LocationType
}

// ReadAll reads the hapefile content.
// If LocationType is InBinary, the file read from binary.
// otherwise reads from file system.
func (c *HapeFile) ReadAll() ([]byte, error) {
	joinedPath := filepath.Join(c.BasePath, c.Name+".hape")
	if c.LocationType == InBinary {
		return Asset(joinedPath)
	}
	return ioutil.ReadFile(joinedPath)
}

// Hapes to get list of hapes
func Hapes() ([]*HapePath, error) {
	hapePaths, err := hapesFromHapePath()
	if err != nil {
		return nil, err
	}
	hapePaths = append(hapePaths, &HapePath{
		Name:         "hapes",
		HapeFiles:    HapesInBinary(),
		LocationType: InBinary,
	})
	return hapePaths, nil
}

func hapesFromHapePath() ([]*HapePath, error) {
	hapePaths := make([]*HapePath, 0)
	hapePath := os.Getenv("HAPEPATH")
	if hapePath == "" {
		return hapePaths, nil
	}
	paths := splitPath(hapePath)
	for _, path := range paths {
		dirEntries, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		path := &HapePath{
			Name:         path,
			HapeFiles:    []string{},
			LocationType: InDirectory,
		}
		for _, entry := range dirEntries {
			name := entry.Name()
			if strings.HasSuffix(name, ".hape") {
				name = strings.TrimSuffix(name, ".hape")
				path.HapeFiles = append(path.HapeFiles, name)
			}
		}
		sort.Strings(path.HapeFiles)
		hapePaths = append(hapePaths, path)
	}
	return hapePaths, nil
}

// GetHape to get hape's ascii art
func (hape *Hape) GetHape() (string, error) {
	src, err := hape.typ.ReadAll()
	if err != nil {
		return "", err
	}

	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\@", "@",
		"\\$", "$",
		"$eyes", hape.eyes,
		"${eyes}", hape.eyes,
		"$tongue", hape.tongue,
		"${tongue}", hape.tongue,
		"$thoughts", string(hape.thoughts),
		"${thoughts}", string(hape.thoughts),
	)
	newsrc := r.Replace(string(src))
	separate := strings.Split(newsrc, "\n")
	mow := make([]string, 0, len(separate))
	for _, line := range separate {
		if strings.Contains(line, "$the_hape = <<EOH") || strings.HasPrefix(line, "##") {
			continue
		}

		if strings.Contains(line, "$ballonOffset = ") {
			line = strings.TrimPrefix(line, "$ballonOffset = ")
			hape.balloonOffset, _ = strconv.Atoi(line)
			continue
		}

		if strings.HasPrefix(line, "EOH") {
			break
		}

		mow = append(mow, line)
	}
	return strings.Join(mow, "\n"), nil
}
