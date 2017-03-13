package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/urfave/cli"
	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
)

func run_cf_example(zone_name string, api_key string, api_email string) {
	// Construct a new API object
	api, err := cloudflare.New(api_key, api_email)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(zone_name) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	api_dns_record_name := fmt.Sprintf("api.%s", zone_name)
	api_dns_record := cloudflare.DNSRecord{Name: api_dns_record_name}
	recs, err := api.DNSRecords(id, api_dns_record)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range recs {
		if !r.Proxied {
			fmt.Printf("%s not proxied. Attempting to update...", r.Content)
			r.Proxied = true
			err := api.UpdateDNSRecord(id, r.ID, r)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("done\n")
		}
	}
}

func merge_config_file(c *cli.Context, config_file *yaml.Yaml) {
	// YELLOW wanna rewrite this so I us a struct but I"m starting to take too much time
	config_flag_slice := []string{"cf_api_key", "cf_api_email", "zone_name"}

	for _, key := range config_flag_slice {
		if c.String(key) == "" && config_file.Get(key) != "" {
			c.Set(key, to.String(config_file.Get(key)))
		}
	}
}

func main() {
	cfpu := cli.NewApp()
	cfpu.Name = "cfpu"
	cfpu.Usage = "Poke the cloudflare dns api and make sure all api entris are proxied."

	cfpu.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load the configuration from `FILE`",
		},
		cli.StringFlag{
			Name:   "cf_api_key",
			Usage:  "Cloudflare API key",
			EnvVar: "CF_API_KEY",
		},
		cli.StringFlag{
			Name:   "cf_api_email",
			Usage:  "Cloudflare API email",
			EnvVar: "CF_API_EMAIL",
		},
		cli.StringFlag{
			Name:  "zone_name, z",
			Usage: "The name of the dns zone you want to interact with",
		},
		cli.StringFlag{
			Name:  "nothing, n",
			Usage: "This is purposefully empty",
		},
	}

	cfpu.Action = func(c *cli.Context) error {
		if c.String("config") != "" {
			settings, err := yaml.Open(c.String("config"))
			if err != nil {
				log.Printf("Could not open YAML file: %s", err.Error())
			}
			merge_config_file(c, settings)
		}
		run_cf_example(c.String("zone_name"), c.String("cf_api_key"), c.String("cf_api_email"))
		return nil
	}

	cfpu.Run(os.Args)
}
