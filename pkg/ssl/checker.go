package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"wzjk-cli/pkg/api"
)

// Checker provides local SSL certificate checking
type Checker struct {
	timeout time.Duration
}

// NewChecker creates a new SSL checker with the specified timeout
func NewChecker(timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Checker{timeout: timeout}
}

// CheckSSL performs a local SSL certificate check for the given domain and port
func (c *Checker) CheckSSL(domain string, port int) (*api.SSLInfo, error) {
	// Parse domain to extract hostname, port, and path
	hostname, parsedPort, path, scheme := parseDomain(domain)

	// Use provided port if specified, otherwise use parsed port
	if port > 0 {
		parsedPort = port
	}

	// Default to 443 if no port specified
	if parsedPort <= 0 {
		parsedPort = 443
	}

	// Determine suggested check type based on scheme
	suggestedCheckType := "https"
	if scheme == "https" {
		suggestedCheckType = "https"
	} else if scheme == "http" {
		suggestedCheckType = "http"
	}

	// Perform TLS connection
	addr := fmt.Sprintf("%s:%d", hostname, parsedPort)
	dialer := &net.Dialer{Timeout: c.timeout}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		ServerName:         hostname,
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, fmt.Errorf("无法连接到 %s: %w", addr, err)
	}
	defer conn.Close()

	// Get peer certificates
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("服务器未提供证书")
	}

	cert := state.PeerCertificates[0]

	// Build SSL info
	info := &api.SSLInfo{
		Domain:             domain,
		Hostname:           hostname,
		Port:               parsedPort,
		Scheme:             scheme,
		Path:               path,
		SuggestedCheckType: suggestedCheckType,
		Issuer:             formatName(cert.Issuer),
		Subject:            formatName(cert.Subject),
		ValidFrom:          cert.NotBefore.Format("2006-01-02 15:04:05"),
		ValidTo:            cert.NotAfter.Format("2006-01-02 15:04:05"),
		DaysUntilExpiry:    daysUntil(cert.NotAfter),
		IsValid:            cert.NotAfter.After(time.Now()) && cert.NotBefore.Before(time.Now()),
		SubjectAltNames:    cert.DNSNames,
	}

	// Check for domain mismatch
	info.DomainMismatch = !matchesDomain(hostname, cert)

	return info, nil
}

// parseDomain extracts hostname, port, path, and scheme from a domain string
func parseDomain(domain string) (hostname string, port int, path string, scheme string) {
	hostname = domain
	scheme = "https"

	// Handle URL format (e.g., "https://example.com:443/path")
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		u, err := url.Parse(domain)
		if err == nil {
			scheme = u.Scheme
			hostname = u.Hostname()
			path = u.Path
			if u.Port() != "" {
				if p, err := strconv.Atoi(u.Port()); err == nil {
					port = p
				}
			}
		}
	} else {
		// Handle hostname with port (e.g., "example.com:443")
		if strings.Contains(hostname, ":") {
			hostParts := strings.Split(hostname, ":")
			hostname = hostParts[0]
			if p, err := strconv.Atoi(hostParts[1]); err == nil && p > 0 && p <= 65535 {
				port = p
			}
		}
	}

	return hostname, port, path, scheme
}

// formatName formats a pkix.Name into a readable string
func formatName(name pkix.Name) string {
	parts := []string{}
	if name.CommonName != "" {
		parts = append(parts, fmt.Sprintf("CN=%s", name.CommonName))
	}
	if len(name.Organization) > 0 {
		parts = append(parts, fmt.Sprintf("O=%s", strings.Join(name.Organization, ", ")))
	}
	if len(name.OrganizationalUnit) > 0 {
		parts = append(parts, fmt.Sprintf("OU=%s", strings.Join(name.OrganizationalUnit, ", ")))
	}
	if len(name.Country) > 0 {
		parts = append(parts, fmt.Sprintf("C=%s", strings.Join(name.Country, ", ")))
	}
	if len(parts) == 0 {
		return name.String()
	}
	return strings.Join(parts, ", ")
}

// daysUntil calculates days until the given time
func daysUntil(t time.Time) int {
	duration := t.Sub(time.Now())
	days := int(duration.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// matchesDomain checks if the certificate matches the given hostname
func matchesDomain(hostname string, cert *x509.Certificate) bool {
	// Check Subject Common Name
	if matchesHostname(hostname, cert.Subject.CommonName) {
		return true
	}

	// Check Subject Alternative Names
	for _, san := range cert.DNSNames {
		if matchesHostname(hostname, san) {
			return true
		}
	}

	return false
}

// matchesHostname checks if a hostname matches a pattern (supports wildcards)
func matchesHostname(hostname, pattern string) bool {
	// Convert both to lowercase for case-insensitive comparison
	hostname = strings.ToLower(hostname)
	pattern = strings.ToLower(pattern)

	// Exact match
	if hostname == pattern {
		return true
	}

	// Wildcard match
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:] // Remove the *
		return strings.HasSuffix(hostname, suffix)
	}

	return false
}
