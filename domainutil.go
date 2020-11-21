package domainutil

import (
	"fmt"
	"regexp"
	"strings"

	ps "golang.org/x/net/publicsuffix"
)

var label = regexp.MustCompile(`[^\.]+`)

func EffectiveTLDPlusN(domain string, n int) (string, error) {
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
