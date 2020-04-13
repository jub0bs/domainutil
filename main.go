package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	ps "golang.org/x/net/publicsuffix"
)

var label = regexp.MustCompile(`[^\.]+`)

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
		host, err := effectiveTLDPlusN(raw, *tldplus)
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

func effectiveTLDPlusN(domain string, n int) (string, error) {
	domain = removeAtMostOneTrailingPeriod(domain)
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return "", fmt.Errorf("empty label in domain %q", domain)
	}
	suffix, _ := ps.PublicSuffix(domain)
	if n == 0 {
		return suffix, nil
	}
	if len(domain) <= len(suffix) { // the domain is only composed of its public suffix
		return "", fmt.Errorf("cannot derive eTLD+%d for domain %q", n, domain)
	}
	i := len(domain) - len(suffix) - 1
	pairs := label.FindAllStringIndex(domain[:i], -1)
	if n > len(pairs) {
		return "", fmt.Errorf("cannot derive eTLD+%d for domain %q", n, domain)
	}
	j := pairs[len(pairs)-n][0]
	return domain[j:], nil
}

func removeAtMostOneTrailingPeriod(s string) string {
	if strings.HasSuffix(s, ".") {
		return s[:len(s)-1]
	}
	return s
}
