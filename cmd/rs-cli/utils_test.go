package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyPrint(t *testing.T) {
	cases := []struct {
		in  int
		out string
	}{
		{2147483647, "Max cash stack, actual price likely higher\n"},
		{1000000000, "1.000.000.000 (1.00b)\n"},
		{1065314566, "1.065.314.566 (1.06b)\n"},
		{1000000, "1.000.000 (1.00m)\n"},
		{1189818, "1.189.818 (1.18m)\n"},
		{1000, "1.000 (1.0k)\n"},
		{14208, "14.208 (14.2k)\n"},
		{100, "100\n"},
	}

	for _, test := range cases {
		out := bytes.Buffer{}
		prettyPrintPrice(&out, test.in)
		assert.Equal(t, test.out, out.String(), fmt.Sprintf("Input was: %d", test.in))
	}
}
