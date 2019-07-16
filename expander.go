package cidr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

//Expand takes a CIDR formatted strings and expands them to corresponding Ipv4 instances including the network and broadcast addresses as first and last elements respectively. Any error returns an empty slice
func Expand(cidr string) []string {

	if strings.Contains(cidr, ":") {
		if ips, ports, err := ExpandWithPort(cidr); err == nil {
			pp := []string{}
			ranges := []string{}
			for _, p := range ports {
				pp = append(pp, fmt.Sprintf("%d", p))
			}
			portsString := strings.Join(pp, ",")

			for _, ip := range ips {
				ranges = append(ranges, fmt.Sprintf("%s:%s", ip, portsString))
			}
			return ranges
		}
		return []string{}
	}

	if !strings.Contains(cidr, "/") {
		//deal with potentially raw IP in non-CIDR format
		cidr += "/32"
	}
	nonCidr := strings.Split(cidr, "/")
	address := nonCidr[0]

	ipAdds, err := net.LookupIP(address)
	if err != nil {
		return []string{}
	}
	combinedIPs := []string{}
	for _, ipAdd := range ipAdds {
		cidr = ipAdd.To4().String() + "/" + nonCidr[1]
		ip, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		size, _ := network.Mask.Size()
		if size == 32 {
			combinedIPs = append(combinedIPs, ip.String())
			continue
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
		combinedIPs = append(combinedIPs, ips...)
	}
	return combinedIPs
}

//ExpandWithPort expands a CIDR range with possible port ranges too
//For example, 10.1.4.2/32:100-200,250,800-855 == 10.1.4.2:100,101,...200,250,800,...,855.
//These will be returned as a slice of IPs, a sorted slice of int ports, and nil if there is no error
func ExpandWithPort(cidrPort string) (cidr []string, ports []int, err error) {
	if split := strings.Split(cidrPort, ":"); len(split) == 2 {
		cidr = Expand(split[0])
		if len(cidr) == 0 {
			return cidr, ports, fmt.Errorf("Invalid CIDR format: %s", split[0])
		}
		if ports, err = expandPorts(split[1]); err != nil {
			return cidr, ports, err
		}

	} else {
		return cidr, ports, fmt.Errorf("Invalid CIDR and port format: %s", cidrPort)
	}
	return
}

func expandPorts(port string) (ports []int, err error) {
	pp := strings.Split(port, ",")
	for _, p := range pp {
		if !strings.Contains(p, "-") {
			intP, e := strconv.Atoi(p)
			if e != nil {
				return ports, e
			}
			ports = append(ports, intP)
		} else if pRange := strings.Split(p, "-"); len(pRange) == 2 {
			low, err := strconv.Atoi(pRange[0])
			if err != nil {
				return ports, fmt.Errorf("Invalid port %s", pRange[0])
			}
			high, err := strconv.Atoi(pRange[1])
			if err != nil {
				return ports, fmt.Errorf("Invalid port %s", pRange[1])
			}
			if low > high {
				low, high = high, low
			}
			for i := low; i <= high; i++ {
				ports = append(ports, i)
			}
		} else {
			return ports, fmt.Errorf("The port range %s is invalid", p)
		}
	}
	return
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
