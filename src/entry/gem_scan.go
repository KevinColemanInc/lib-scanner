/*
 */

package main

import (
	"github.com/KevinColemanInc/lib-crawl/src/handle"

	"github.com/tucnak/climax"
	"fmt"
)

func main() {
	fmt.Println("Started!")
	demo := climax.New("gem_scan")
	demo.Brief = "Takes in a directory to scan all of the gems for unexpected data access."
	demo.Version = "0.0.1"

	joinCmd := climax.Command{
		Name:  "scan",
		Brief: "scan directory for suspicious ruby code",
		Usage: ``,
		Help:  ``,

		Examples: []climax.Example{
			{
				Usecase:     `~/.rvm/gems/ruby-2.7.0/`,
				Description: `Scans all ruby gems using installed with 2.7`,
			},
		},

		Handle: handle.Scan}

	demo.AddCommand(joinCmd)
	demo.Run()
}
