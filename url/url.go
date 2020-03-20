package url

import (
	"errors"
	"strings"
)

// Type represents the type of an URL
type Type uint8

const (
	// Absolute is the type of an URL which is absolute with respect to the root
	Absolute Type = iota

	// Relative is the type of an URL which is relative with respect to the calling URL
	Relative

	// External is the type of an URL which points to a resource in a different domain
	External

	// RedundantlyVerbose is the type of an URL which uses a FQDN to refer to a resource in the local domain
	RedundantlyVerbose

	// Anchor is a same-page anchor link
	Anchor

	// Mail is a mailto link
	Mail
)

func (t Type) String() string {
	switch t {
	case Absolute:
		return "Absolute"
	case Relative:
		return "Relative"
	case External:
		return "External"
	case RedundantlyVerbose:
		return "Redundantly Verbose"
	case Anchor:
		return "Anchor"
	}

	panic(errors.New("invalid Type value"))
}

// Deducer deduces the type of a URL
type Deducer struct {
	LocalDomain      string
	localDomainSet   bool
	localDomainHTTP  string
	localDomainHTTPS string
}

// NewDeducer creates a new Deducer that cannot detect RedundantlyVerbose URLs
func NewDeducer() *Deducer {
	return &Deducer{
		localDomainSet: false,
	}
}

// NewDeducerWithLocalDomain creates a new Deducer that can detect RedundantlyVerbose URLs
func NewDeducerWithLocalDomain(localDomain string) *Deducer {
	return &Deducer{
		LocalDomain:      localDomain,
		localDomainHTTP:  "http://" + localDomain,
		localDomainHTTPS: "https://" + localDomain,
		localDomainSet:   true,
	}
}

func (*Deducer) ContainsAmpersand(dest string) bool {
	return strings.Contains(dest, "&")
}

// DeduceTypeOfDestination determines the Type of the given URL string
func (d *Deducer) DeduceTypeOfDestination(dest string) Type {
	if strings.HasPrefix(dest, "#") {
		return Anchor
	}

	if strings.HasPrefix(dest, "mailto:") {
		return Mail
	}

	if strings.HasPrefix(dest, "/") {
		return Absolute
	}

	if isAbsolute(dest) {
		if d.localDomainSet && isRedundantFullURL(dest, d.localDomainHTTP, d.localDomainHTTPS) {
			return RedundantlyVerbose
		}

		return External
	}

	return Relative
}

func (d *Deducer) RewriteRedundantlyVerboseLink(dest string) string {
	if strings.HasPrefix(dest, d.localDomainHTTP) {
		return strings.Replace(dest, d.localDomainHTTP, "", 1)
	}

	if strings.HasPrefix(dest, d.localDomainHTTPS) {
		return strings.Replace(dest, d.localDomainHTTPS, "", 1)
	}

	return dest
}

func (*Deducer) RewriteRelativeLink(dest, filePath string) string {
	//fmt.Printf("Processing links in file: %q\n", filePath)
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}

	//fmt.Printf("Processing links in file: %q\n", filePath)
	filePath = stripFileName(filePath)

	return filePath + "/" + dest
}

func (*Deducer) RewriteContainsAmpersandLink(dest string) string {
	return strings.ReplaceAll(dest, "&", "")
}

func isAbsolute(dest string) bool {
	if strings.Contains(dest, "://") {
		return true
	}

	return false
}

func isRedundantFullURL(url, localDomainHTTP, localDomainHTTPS string) bool {
	if strings.HasPrefix(url, localDomainHTTP) {
		return true
	}

	if strings.HasPrefix(url, localDomainHTTPS) {
		return true
	}

	return false
}

func stripFileName(filePath string) string {
	if index := strings.LastIndex(filePath, "/"); index != -1 {
		return filePath[:index]
	}

	return filePath
}
