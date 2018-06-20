# CIDR 
CIDR is a simple utility to generate the IPv4 addresses in a CIDR range. It could also be used to check the membership of an IP v4 address in a CIDR range.

## Using as a library
In order to start, go get this repository:
```go
go get  github.com/adedayo/net/cidr
```

### Example
In your code simply import as usual and enjoy:

```go
package main

import "github.com/adedayo/net/cidr"

func main() {
	ips := cidr.Expand("192.168.2.5/30")

	for _, ip := range ips {
		println(ip)
	}
}

```
This should print the set of IPs described by the CIDR range 192.168.2.5/30:
```
192.168.2.4
192.168.2.5
192.168.2.6
192.168.2.7
```

## Using it as a command-line tool
CIDR is also available as a command-line tool
### Generating IPs in a CIDR range

```bash
cidr 192.168.2.5/30 10.11.12.13/31
```

This should generate a simply formatted output:

```bash
192.168.2.5/30: 192.168.2.4 192.168.2.5 192.168.2.6 192.168.2.7

10.11.12.13/31: 10.11.12.12 10.11.12.13
```

For a JSON-formatted result, use the JSON `-j` or `--json` flag:

```bash
cidr --json 192.168.2.5/30 10.11.12.13/31
```
This should produce:

```json
{
"192.168.2.5/30": ["192.168.2.4", "192.168.2.5", "192.168.2.6", "192.168.2.7"],

"10.11.12.13/31": ["10.11.12.12", "10.11.12.13"]
}
```

### Checking CIDR range membership

The structure of the checking command is as follows 
```bash
cidr [flag] check <space separated list of CIDR ranges> <delimiter> <space-separated list of IP addresses to check>
```

The delimiter can be any of `contains`, `c` or simply `,`

Examples

To check the membership of the IP addresses 192.168.10.3 and 192.168.10.9 in the CIDR ranges 192.168.10.1/30 and 192.168.10.1/28 run
```bash
cidr check 192.168.10.1/30 192.168.10.1/28 contains 192.168.10.3 192.168.10.9
```

The result is 
```
192.168.10.1/30: 192.168.10.3,true 192.168.10.9,false
192.168.10.1/28: 192.168.10.3,true 192.168.10.9,true
```

For a JSON output the above is equivalent to

```bash
cidr --json check 192.168.10.1/30 192.168.10.1/28 , 192.168.10.3 192.168.10.9
```

This produces
```json
{
"192.168.10.1/30": [{"ip":"192.168.10.3","belongs":true},{"ip":"192.168.10.9","belongs":false}],
"192.168.10.1/28": [{"ip":"192.168.10.3","belongs":true},{"ip":"192.168.10.9","belongs":true}]
}
```

## License
BSD 3-Clause License