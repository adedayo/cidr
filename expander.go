package cidr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

//Expand takes a CIDR formatted strings and expands them to corresponding Ipv4 instances including the network and broadcast addresses as first and last elements respectively. Any error returns an empty slice
func Expand(cidr string) []string {
	if !strings.Contains(cidr, "/") {
		//deal with potentially raw IP in non-CIDR format
		cidr += "/32"
	}
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return []string{}
	}
	size, _ := network.Mask.Size()
	if size == 32 {
		return []string{ip.String()}
	}

	//include network and broadcast addresses in count
	hostCount := 2 << uint(31-size)
	ips := make([]string, hostCount)
	ips[0] = network.IP.String()

	//network address starting point
	octets := decimalOctets(network.IP)

	for i := 1; i < hostCount; i++ {
		octets[3]++
		if octets[3] > 255 {
			octets[3] = 0
			octets[2]++
		}

		if octets[2] > 255 {
			octets[2] = 0
			octets[1]++
		}

		if octets[1] >= 255 {
			octets[1] = 0
			octets[0]++
		}
		ips[i] = toIP(octets)
	}
	return ips
}

func toIP(oct []int) string {
	if len(oct) != 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s.%s", strconv.Itoa(oct[0]), strconv.Itoa(oct[1]),
		strconv.Itoa(oct[2]), strconv.Itoa(oct[3]))
}

func decimalOctets(ip net.IP) []int {
	result := make([]int, 4)
	octets := ip.To4()
	for i := 0; i < 4; i++ {
		result[i] = int(octets[i])
	}
	return result
}

//Membership describes membership status of an IP in a CIDR range
type Membership struct {
	CIDR    string
	IP      string
	Belongs bool
}

//Contains returns the membership of a set of IPs in a CIDR range
func Contains(cidr string, ips ...string) (mem []Membership) {
	//IPNet seems to exclude the last (broadcast address in the Contains() method, so using my own Expand() for membership test)
	ipAdds := Expand(cidr)
	ipMap := make(map[string]bool, len(ipAdds))
	for _, ip := range ipAdds {
		ipMap[ip] = true
	}
	for _, ip := range ips {
		_, belongs := ipMap[ip]
		mem = append(mem, Membership{
			CIDR:    cidr,
			IP:      ip,
			Belongs: belongs,
		})
	}
	return
}
