package main

import (
	"fmt"
	"net/url"
)

func clean(host, href string) (parsedHrefRet string, err error) {
	// Make host string into URL struct
	hostURL, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	// Make href into URL struct
	parsedHref, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	// Case for path fragment -> #foo
	// Use host path insetad of fragment
	// The main URL should be crawled or de-duped by the crawler
	if len(parsedHref.Path) == 0 && len(parsedHref.Fragment) != 0 {
		parsedHref.Path = hostURL.Path
		parsedHref.Fragment = ""
	}

	// Assign href sheme to host scheme, if empty
	if parsedHref.Scheme == "" {
		parsedHref.Scheme = hostURL.Scheme
	} else if parsedHref.Scheme != "https" && parsedHref.Scheme != "http" {
		return "", fmt.Errorf("%s scheme is invalid", parsedHref.Scheme)
	}

	// Assign href hostname to host hostname, if empty
	if parsedHref.Host == "" {
		parsedHref.Host = hostURL.Host
	} else if parsedHref.Host != hostURL.Host {
		return "", fmt.Errorf("%s hostname is invalid", parsedHref.Host)
	}

	// Return href as string to empty error message
	return parsedHref.String(), nil
}
