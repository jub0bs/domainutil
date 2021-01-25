package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jub0bs/domainutil"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, w io.Writer, r io.Reader) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	tldplus := flags.Int("tldplus", 1, "level above public suffix")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Text()
		host, _, err := domainutil.EffectiveTLDPlusN(raw, *tldplus)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, host)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading standard input: %v", err)
	}
	return nil
}
