package osc

import (
	"regexp"
	"strings"
)

/*
MessageHandleFunc is a function type that accepts a pointer to a Message.
*/
type MessageHandleFunc func(*Message)

/*
Method represents an address pattern with associated invokable function.
*/
type Method struct {
	AddressPattern string
	Function       MessageHandleFunc
	regexp         *regexp.Regexp
}

/*
AddressSpace holds a set of methods that an OSC server can respond to.
*/
type AddressSpace struct {
	methods []Method
}

/*
Handle adds an OSC method to the AddressSpace. If the AddressPattern is of invalid format, an error is returned.
*/
func (a *AddressSpace) Handle(addressPattern string, fn MessageHandleFunc) error {
	// Compile a regexp to use when matching the address pattern
	regexp, err := addressPatternToRegexp(addressPattern)
	if err != nil {
		return err
	}

	method := Method{
		AddressPattern: addressPattern,
		Function:       fn,
		regexp:         regexp,
	}

	a.methods = append(a.methods, method)

	return nil
}

/*
Methods returns the OSC methods held in an AddressSpace.
*/
func (a AddressSpace) Methods() []Method {
	return a.methods
}

/*
Dispatch finds a matching OSC method for the Message m, and invokes it if found.
*/
func (a AddressSpace) Dispatch(m *Message) {
	if m == nil {
		return
	}

	for _, h := range a.methods {
		if h.regexp.MatchString(m.Address) {
			h.Function(m)
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
