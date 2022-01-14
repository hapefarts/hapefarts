package hapesay

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHape_Clone(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		from *Hape
		want *Hape
	}{
		{
			name: "without options",
			opts: []Option{},
			from: func() *Hape {
				hape, _ := New()
				return hape
			}(),
			want: func() *Hape {
				hape, _ := New()
				return hape
			}(),
		},
		{
			name: "with some options",
			opts: []Option{},
			from: func() *Hape {
				hape, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return hape
			}(),
			want: func() *Hape {
				hape, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return hape
			}(),
		},
		{
			name: "clone and some options",
			opts: []Option{
				Thinking(),
				Thoughts1('o'),
				Thoughts2('o'),
			},
			from: func() *Hape {
				hape, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return hape
			}(),
			want: func() *Hape {
				hape, _ := New(
					Type("docker"),
					BallonWidth(60),
					Thinking(),
					Thoughts1('o'),
					Thoughts2('o'),
				)
				return hape
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.want.Clone(tt.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(Hape{}),
				cmpopts.IgnoreFields(Hape{}, "buf")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}

	t.Run("random", func(t *testing.T) {
		hape, _ := New(
			Type(""),
			Thinking(),
			Thoughts1('o'),
			Thoughts2('o'),
			Eyes("xx"),
			Tongue("u"),
			Random(),
		)

		cloned, _ := hape.Clone()

		if diff := cmp.Diff(hape, cloned,
			cmp.AllowUnexported(Hape{}),
			cmpopts.IgnoreFields(Hape{}, "buf")); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("error", func(t *testing.T) {
		hape, err := New()
		if err != nil {
			t.Fatal(err)
		}

		wantErr := errors.New("error")
		_, err = hape.Clone(func(*Hape) error {
			return wantErr
		})
		if wantErr != err {
			t.Fatalf("want %v, but got %v", wantErr, err)
		}
	})
}

func Test_adjustTo2Chars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "  ",
		},
		{
			name: "1 character",
			s:    "1",
			want: "1 ",
		},
		{
			name: "2 characters",
			s:    "12",
			want: "12",
		},
		{
			name: "3 characters",
			s:    "123",
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := adjustTo2Chars(tt.s); got != tt.want {
				t.Errorf("adjustTo2Chars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFound_Error(t *testing.T) {
	file := "test"
	n := &NotFound{
		Hapefile: file,
	}
	want := fmt.Sprintf("not found %q hapefile", file)
	if want != n.Error() {
		t.Fatalf("want %q but got %q", want, n.Error())
	}
}
