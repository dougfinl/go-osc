package osc

import (
	"regexp"
	"strings"
)

type messageHandler struct {
	addressPattern string
	regexp         *regexp.Regexp
	function       func(m *Message)
}

type messageDispatcher struct {
	handlers []messageHandler
}

func (d *messageDispatcher) addHandler(addressPattern string, function func(*Message)) error {
	// Compile a regexp to use when matching the address pattern
	regexp, err := addressPatternToRegexp(addressPattern)
	if err != nil {
		return err
	}

	handler := messageHandler{addressPattern: addressPattern, regexp: regexp, function: function}
	d.handlers = append(d.handlers, handler)

	return nil
}

func (d messageDispatcher) dispatch(m *Message) {
	if m == nil {
		return
	}

	for _, h := range d.handlers {
		if h.regexp.MatchString(m.Address) {
			h.function(m)
		}
	}
}

/*
addressPatternToRegexp creates a regular expression used to efficiently match the address pattern.
*/
func addressPatternToRegexp(addressPattern string) (*regexp.Regexp, error) {
	// Escape forward slashes
	apRegexp := strings.Replace(addressPattern, "/", "\\/", -1)
	// apRegexp := addressPattern

	// Convert basic wildcard expressions and classes
	apRegexp = strings.Replace(apRegexp, "?", ".", -1)
	apRegexp = strings.Replace(apRegexp, "*", ".*", -1)
	apRegexp = strings.Replace(apRegexp, "![", "[^", -1)

	// Convert group notation
	apRegexp = strings.Replace(apRegexp, "{", "(", -1)
	apRegexp = strings.Replace(apRegexp, "}", ")", -1)
	apRegexp = strings.Replace(apRegexp, ",", "|", -1)

	re, err := regexp.Compile(apRegexp)
	if err != nil {
		return nil, err
	}

	return re, nil
}

/*
match will return true if address is a valid match for addressPattern.
*/
func match(address string, addressPattern string) bool {
	return false
}
