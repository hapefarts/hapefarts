package hapesay

import (
	"embed"
	"sort"
	"strings"
)

//go:embed hapes/*
var hapesDir embed.FS

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(path string) ([]byte, error) {
	return hapesDir.ReadFile(path)
}

// AssetNames returns the list of filename of the assets.
func AssetNames() []string {
	entries, err := hapesDir.ReadDir("hapes")
	if err != nil {
		panic(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".hape")
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

var hapesInBinary = AssetNames()

// HapesInBinary returns the list of hapefiles which are in binary.
// the list is memoized.
func HapesInBinary() []string {
	return hapesInBinary
}
