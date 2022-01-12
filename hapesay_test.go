package hapesay

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHapes(t *testing.T) {
	t.Run("no set COWPATH env", func(t *testing.T) {
		hapePaths, err := Hapes()
		if err != nil {
			t.Fatal(err)
		}
		if len(hapePaths) != 1 {
			t.Fatalf("want 1, but got %d", len(hapePaths))
		}
		hapePath := hapePaths[0]
		if len(hapePath.HapeFiles) == 0 {
			t.Fatalf("no hapefiles")
		}

		wantHapePath := &HapePath{
			Name:         "hapes",
			LocationType: InBinary,
		}
		if diff := cmp.Diff(wantHapePath, hapePath,
			cmpopts.IgnoreFields(HapePath{}, "HapeFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("set COWPATH env", func(t *testing.T) {
		hapepath := filepath.Join("testdata", "testdir")

		os.Setenv("COWPATH", hapepath)
		defer os.Unsetenv("COWPATH")

		hapePaths, err := Hapes()
		if err != nil {
			t.Fatal(err)
		}
		if len(hapePaths) != 2 {
			t.Fatalf("want 2, but got %d", len(hapePaths))
		}

		wants := []*HapePath{
			{
				Name:         "testdata/testdir",
				LocationType: InDirectory,
			},
			{
				Name:         "hapes",
				LocationType: InBinary,
			},
		}
		if diff := cmp.Diff(wants, hapePaths,
			cmpopts.IgnoreFields(HapePath{}, "HapeFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}

		if len(hapePaths[0].HapeFiles) != 1 {
			t.Fatalf("unexpected hapefiles len = %d, %+v",
				len(hapePaths[0].HapeFiles), hapePaths[0].HapeFiles,
			)
		}

		if hapePaths[0].HapeFiles[0] != "test" {
			t.Fatalf("want %q but got %q", "test", hapePaths[0].HapeFiles[0])
		}
	})

	t.Run("set COWPATH env", func(t *testing.T) {
		os.Setenv("COWPATH", "notfound")
		defer os.Unsetenv("COWPATH")

		_, err := Hapes()
		if err == nil {
			t.Fatal("want error")
		}
	})

}

func TestHapePath_Lookup(t *testing.T) {
	t.Run("looked for hapefile", func(t *testing.T) {
		c := &HapePath{
			Name:         "basepath",
			HapeFiles:    []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("test")
		if !ok {
			t.Errorf("want %v", ok)
		}
		want := &HapeFile{
			Name:         "test",
			BasePath:     "basepath",
			LocationType: InBinary,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("no hapefile", func(t *testing.T) {
		c := &HapePath{
			Name:         "basepath",
			HapeFiles:    []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("no hapefile")
		if ok {
			t.Errorf("want %v", !ok)
		}
		if got != nil {
			t.Error("want nil")
		}
	})
}

func TestHapeFile_ReadAll(t *testing.T) {
	fromTestData := &HapeFile{
		Name:         "test",
		BasePath:     filepath.Join("testdata", "testdir"),
		LocationType: InDirectory,
	}
	fromTestdataContent, err := fromTestData.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	fromBinary := &HapeFile{
		Name:         "mobile",
		BasePath:     "hapes",
		LocationType: InBinary,
	}
	fromBinaryContent, err := fromBinary.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fromTestdataContent, fromBinaryContent) {
		t.Fatalf("testdata\n%s\n\nbinary%s\n", string(fromTestdataContent), string(fromBinaryContent))
	}

}

const defaultSay = ` ________ 
< hapesay >
 -------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`

func TestSay(t *testing.T) {
	type args struct {
		phrase  string
		options []Option
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "mobile",
			args: args{
				phrase: "hello!",
			},
			wantFile: "mobile.hape",
			wantErr:  false,
		},
		{
			name: "nest",
			args: args{
				phrase: defaultSay,
				options: []Option{
					DisableWordWrap(),
				},
			},
			wantFile: "nest.hape",
			wantErr:  false,
		},
		{
			name: "error",
			args: args{
				phrase: "error",
				options: []Option{
					func(*Hape) error {
						return errors.New("error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Say(tt.args.phrase, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Say() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			filename := filepath.Join("testdata", tt.wantFile)
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			want := string(content)
			if want != got {
				t.Fatalf("want\n%s\n\ngot\n%s", want, got)
			}
		})
	}
}
