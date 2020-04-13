package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestHappyPaths(t *testing.T) {
	cases := []struct {
		args []string
		in   string
		want string
	}{
		{
			[]string{"main", "--tldplus=-1"}, // parses that as 1
			"",
			"",
		},
		{
			[]string{"main", "--tldplus=1"},
			"foo.com\nbar.co.uk",
			"foo.com\nbar.co.uk\n",
		},
		{
			[]string{"main"},
			"foo.com\nbar.co.uk",
			"foo.com\nbar.co.uk\n",
		},
		{
			[]string{"main"},
			"foo.com.\nbar.co.uk.",
			"foo.com\nbar.co.uk\n",
		},
		{
			[]string{"main", "--tldplus=0"},
			"foo.com\nbar.co.uk",
			"com\nco.uk\n",
		},
		{
			[]string{"main", "--tldplus=2"},
			"admin.foo.com\nadmin.bar.co.uk",
			"admin.foo.com\nadmin.bar.co.uk\n",
		},
	}

	for _, c := range cases {
		r := newReader(c.in)
		var w strings.Builder
		err := run(c.args, &w, r)
		got := w.String()
		if err != nil || got != c.want {
			t.Errorf("want: %q, got: %q", c.want, got)
		}
	}
}

func TestUnhappyPaths(t *testing.T) {
	cases := []struct {
		args []string
		in   string
	}{
		{
			[]string{"main"},
			"foo.com..\nbar.co.uk..",
		},
		{
			[]string{"main"},
			".foo.com\n.bar.co.uk",
		},
		{
			[]string{"main"},
			"foo..com\nbar..co.uk",
		},
		{
			[]string{"main"},
			"com\nco.uk",
		},
		{
			[]string{"main", "--tldplus=2"},
			"foo.com\nbar.co.uk",
		},
		{
			[]string{"main", "--fantasyflag"},
			"foo.com\nbar.co.uk",
		},
	}

	for _, c := range cases {
		r := newReader(c.in)
		var w strings.Builder
		err := run(c.args, &w, r)
		got := w.String()
		if err == nil || got != "" {
			t.Errorf("want: \"\", got: %q", got)
		}
	}
}

func newReader(ss ...string) io.Reader {
	const sep = "\n"
	s := strings.Join(ss, sep)
	return strings.NewReader(s)
}

func TestScannerFailure(t *testing.T) {
	r := &failingReader{}
	var w strings.Builder
	err := run([]string{"main"}, &w, r)
	want := ""
	got := w.String()
	if err == nil || got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}

type failingReader struct{}

func (*failingReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("some I/O error")
}
