package main

import (
	"fmt"
	"github.com/adedayo/cidr"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

var (
	//Version cidr version
	Version = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "cidr"
	app.Version = Version
	app.Usage = "Expand CIDR range to individual IP addresses"
	app.UsageText = `Expands a CIDR range or a space separated list of CIDR ranges to individual IP addresses. 
	
Examples:
	
cidr expand 8.8.8.8/24 192.168.10.1/30

or simply

cidr 8.8.8.8/24 192.168.10.1/30

`
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "generate JSON output",
		},
	}
	app.EnableBashCompletion = true

	app.Authors = []cli.Author{
		{
			Name:  "Adedayo Adetoye (Dayo)",
			Email: "https://github.com/adedayo",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "expand",
			Aliases: []string{"e"},
			Usage:   "Expand one or more space-separated `CIDR ranges` into IP addresses",
			UsageText: `
	cidr 8.8.8.8/24 192.168.10.1/30
	
	To generate a JSON output, pass the -j or --json flag:
	cidr --json 8.8.8.8/24 192.168.10.1/30
			`,
			Action: func(c *cli.Context) error {
				return process(c)
			},
		},
		cli.Command{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Checks whether one or more space-separated `CIDR ranges` contain one or more space-separated `IP addresses`",
			UsageText: `
	Check a single IP in a CIDR range:
	cidr check 192.168.10.1/30 contains 192.168.10.3

	To check multiple CIDRs and IPs, generating a JSON output:
	cidr check 192.168.10.1/30 220.10.5.15/28 contains 192.168.10.3 220.10.5.18 192.168.10.230 			
			`,
			Action: func(c *cli.Context) error {
				return check(c)
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		return process(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func process(c *cli.Context) error {

	if c.NArg() == 0 {
		c.App.Run([]string{"cidr", "h"})
	} else {
		if c.GlobalBool("json") {
			result := "{\n"
			prefix := ""
			for i := 0; i < c.NArg(); i++ {
				arg := c.Args().Get(i)
				if i > 0 {
					prefix = ",\n\n"
				}
				ips := cidr.Expand(arg)
				for ind, ip := range ips {
					ips[ind] = fmt.Sprintf(`"%s"`, ip)
				}
				result += fmt.Sprintf("%s\"%s\": [%s]", prefix, arg, strings.Join(ips, ", "))
			}
			result += "\n}"
			fmt.Println(result)

		} else {
			for i := 0; i < c.NArg(); i++ {
				arg := c.Args().Get(i)
				fmt.Printf("%s: %s\n\n", arg, strings.Join(cidr.Expand(arg), " "))
			}
		}

	}
	return nil
}

func check(c *cli.Context) error {
	if c.NArg() < 3 {
		c.App.Run([]string{"cidr", "help", "check"})
	} else {
		cidrs := []string{}
		ips := []string{}
		isCider := true
		for i := 0; i < c.NArg(); i++ {
			token := c.Args().Get(i)
			if isCider && (token == "contains" || token == "," || token == "c") {
				isCider = false
				continue
			}
			if isCider {
				cidrs = append(cidrs, token)
			} else {
				ips = append(ips, token)
			}
		}
		result := ""
		json := c.GlobalBool("json")
		if json {
			result += "{\n"
			firstCidr := true
			for _, cid := range cidrs {
				mem := cidr.Contains(cid, ips...)

				prefixCidr := ",\n"
				if firstCidr {
					prefixCidr = ""
					firstCidr = false
				}

				result += fmt.Sprintf(`%s"%s": [`, prefixCidr, cid)
				first := true
				for _, m := range mem {
					prefix := ","
					if first {
						prefix = ""
						first = false
					}
					result += fmt.Sprintf(`%s{"ip":"%s","belongs":%t}`, prefix, m.IP, m.Belongs)
				}
				result += "]"

			}
			result += "\n}"
		} else {
			result = ""
			firstCidr := true
			for _, cid := range cidrs {
				mem := cidr.Contains(cid, ips...)
				prefixCidr := "\n"
				if firstCidr {
					prefixCidr = ""
					firstCidr = false
				}
				result += fmt.Sprintf("%s%s: ", prefixCidr, cid)
				first := true
				for _, m := range mem {
					prefix := " "
					if first {
						prefix = ""
						first = false
					}
					result += fmt.Sprintf(`%s%s,%t`, prefix, m.IP, m.Belongs)
				}
			}
		}

		println(result)
	}
	return nil
}
