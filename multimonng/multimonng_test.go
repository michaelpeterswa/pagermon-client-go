package multimonng_test

import (
	"testing"

	"github.com/michaelpeterswa/go-lib/multimonng"
	"github.com/stretchr/testify/assert"
)

func TestParseMultimonLine(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected *multimonng.MultimonNGMessage
	}{
		{
			Name:  "Simple Parse",
			Input: "POCSAG1200: Address: 1234567  Function: 0  Alpha:   Aid - Emergency; *FTAC - 1*;  Test Emergency Location; 7xxx Test Rd NE, RM; A1; 47.0;-122.0<EOT><NUL>",
			Expected: &multimonng.MultimonNGMessage{
				Mode:     "POCSAG1200",
				Address:  "1234567",
				Function: "0",
				Alpha:    "Aid - Emergency; *FTAC - 1*;  Test Emergency Location; 7xxx Test Rd NE, RM; A1; 47.0;-122.0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := multimonng.ParseMultimonLine(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, res)
		})
	}
}
