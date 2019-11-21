package main

import (
	"time"

	"github.com/cloudflare/cfssl/config"
)

var defaultConfig = config.DefaultConfig()

// ServerConfig is taken from https://github.com/cloudflare/cfssl/blob/master/cli/printdefault/defaults.go
var ServerConfig = config.Config{
	Signing: &config.Signing{
		Default: &config.SigningProfile{
			Expiry:       168 * time.Hour,
			ExpiryString: "168h",
		},
		Profiles: map[string]*config.SigningProfile{
			"www": &config.SigningProfile{
				Expiry:       defaultConfig.Expiry,
				ExpiryString: defaultConfig.ExpiryString,
				Usage:        []string{"signing", "key encipherment", "server auth"},
			},
		},
	},
}
