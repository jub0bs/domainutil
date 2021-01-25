package domainutil

import (
	"fmt"
	"regexp"
	"strings"

	ps "golang.org/x/net/publicsuffix"
)

const pkgName = "domainutil"

var label = regexp.MustCompile(`[^\.]+`)

func EffectiveTLDPlusN(domain string, n int) (string, bool, error) {
	domain = removeAtMostOneTrailingPeriod(domain)
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return "", false, fmt.Errorf("%s: empty label in domain %q", pkgName, domain)
	}
	suffix, icann := ps.PublicSuffix(domain)
	if n == 0 {
		return suffix, icann, nil
	}
	if len(domain) <= len(suffix) { // the domain is only composed of its public suffix
		return "", false, fmt.Errorf("%s: cannot derive eTLD+%d for domain %q", pkgName, n, domain)
	}
	i := len(domain) - len(suffix) - 1
	pairs := label.FindAllStringIndex(domain[:i], -1)
	if n > len(pairs) {
		return "", false, fmt.Errorf("%s: cannot derive eTLD+%d for domain %q", pkgName, n, domain)
	}
	j := pairs[len(pairs)-n][0]
	return domain[j:], icann, nil
}

func removeAtMostOneTrailingPeriod(s string) string {
	return strings.TrimSuffix(s, ".")
}
