package cidr

import (
	"fmt"
	"net"
	"testing"
)

var (
	baseIP = "192.168.0.2"
)

func TestExpand(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{
			input:    32,
			expected: 1,
		},
		{
			input:    24,
			expected: 256,
		},
		{
			input:    16,
			expected: 65536,
		},
		{
			input:    8,
			expected: 16777216,
		},
	}
	for _, tc := range testCases {
		t.Run("Expected result lengths", func(t *testing.T) {
			obtained := len(Expand(fmt.Sprintf("%s/%d", baseIP, tc.input)))
			if obtained != tc.expected {
				t.Errorf("Number of IPs in %s/%d should be %d, but got %d", baseIP, tc.input, tc.expected, obtained)
			}
		})
	}
}
func TestBadArgument(t *testing.T) {
	cidr := fmt.Sprintf("%s/", baseIP)
	expanded := Expand(cidr)
	if len(expanded) != 0 {
		t.Error("Bad argument should return an empty slice")
	}
}

func TestOctetConversion(t *testing.T) {
	ip := net.ParseIP(baseIP)
	dec := []int{192, 168, 0, 2}
	oct := decimalOctets(ip)
	for i, x := range dec {
		t.Run("Conversion of IP to octets", func(t *testing.T) {
			if x != oct[i] {
				t.Errorf("Invalid conversion of %s to decimal octets should be %#v, but got index %d as %d instead of %d", baseIP, dec, i, oct[i], x)
			}
		})
	}
}

func TestIPConversion(t *testing.T) {
	dec := []int{192, 168, 0, 2}
	if toIP(dec) != baseIP {
		t.Errorf("Invalid conversion of decimal octets %#v to string, should be %s", dec, baseIP)
	}

	if toIP(dec[0:2]) != "" {
		t.Error("Bad octet ought to return empty IP")
	}
}

func TestExpandBadCIDR(t *testing.T) {
	exp := Expand("255.256.256.255/24")
	if len(exp) != 0 {
		t.Errorf("Bad CIDR should return empty list, instead of %#v", exp)
	}
}

func TestMembership(t *testing.T) {
	membership := Contains("10.10.10.3/30", "10.10.10.0", "10.10.10.1", "10.10.10.2", "10.10.10.3")
	for _, member := range membership {
		if !member.Belongs {
			t.Errorf("The IP %s should belong to the CIDR range %s", member.IP, member.CIDR)
		}
	}
}
