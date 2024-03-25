// Package multimonng provides parsing tools to convert multimon-ng output lines into a format that can be sent to a PagerMon server
package multimonng

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MultimonNGMessage represents a message that has been parsed from multimon-ng output
type MultimonNGMessage struct {
	Mode     string // Mode is the decoding mode of the message
	Address  string // Address is the capcode of the message
	Function string // Function is the function of the message
	Alpha    string // Alpha is the alphanumeric (not guaranteed) message content
}

var (
	addressRegex  = regexp.MustCompile(`Address: (\d{1,7})`)  // Should only be able to be between 1 and 7 digits
	functionRegex = regexp.MustCompile(`Function: (\d{1,2})`) // Unknown if this is correct or what the possible values are
	alphaRegex    = regexp.MustCompile(`Alpha: (.*)`)         // Alpha is the alphanumeric (not guaranteed) message content

	ErrAddressNotFound  = fmt.Errorf("address not found")  // Error message for when address is not found
	ErrFunctionNotFound = fmt.Errorf("function not found") // Error message for when function is not found
	ErrAlphaNotFound    = fmt.Errorf("alpha not found")    // Error message for when alpha is not found
)

// newMultimonNGMessage creates a new MultimonNGMessage
func newMultimonNGMessage(mode string, address string, function string, alpha string) *MultimonNGMessage {
	return &MultimonNGMessage{Mode: mode, Address: address, Function: function, Alpha: alpha}
}

// ParseMultimonLine parses a line of multimon-ng output and returns a MultimonNGMessage or an error
func ParseMultimonLine(line string) (*MultimonNGMessage, error) {
	mode := strings.TrimSuffix(strings.SplitAfter(line, ":")[0], ":")

	address := "0"
	addressMatches := addressRegex.FindStringSubmatch(line)
	if len(addressMatches) > 1 {
		i, err := strconv.ParseInt(addressMatches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing address: %w", err)
		}
		address = fmt.Sprintf("%07d", i)
	} else {
		return nil, fmt.Errorf("error parsing address: %w", ErrAddressNotFound)
	}

	function := "0"
	functionMatches := functionRegex.FindStringSubmatch(line)
	if len(functionMatches) > 1 {
		function = functionMatches[1]
	} else {
		return nil, fmt.Errorf("error parsing function: %w", ErrFunctionNotFound)
	}

	alpha := ""
	alphaMatches := alphaRegex.FindStringSubmatch(line)
	if len(alphaMatches) > 1 {
		alpha = trimEndSequences(strings.TrimLeft(alphaMatches[1], " ")) // Trim all left spaces and remove <EOT> and <NUL> sequences
	} else {
		return nil, fmt.Errorf("error parsing alpha: %w", ErrAlphaNotFound)
	}

	return newMultimonNGMessage(mode, address, function, alpha), nil
}

// trimEndSequences trims the <EOT> and <NUL> sequences from the end of a string
func trimEndSequences(input string) string {
	for strings.HasSuffix(input, "<EOT>") || strings.HasSuffix(input, "<NUL>") {
		input = strings.TrimSuffix(input, "<EOT>")
		input = strings.TrimSuffix(input, "<NUL>")
	}
	return input
}
