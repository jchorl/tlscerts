package main

import (
	"time"

	"github.com/cloudflare/cfssl/config"
)

var defaultConfig = config.DefaultConfig()

// DefaultServerConfig is the default signing config for server certs.
// It takes defaults from https://github.com/cloudflare/cfssl/blob/master/cli/printdefault/defaults.go
func DefaultServerConfig() config.Config {
	return config.Config{
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
}

// DefaultClientConfig is the default signing config for client certs.
// It takes defaults from https://github.com/cloudflare/cfssl/blob/master/cli/printdefault/defaults.go
func DefaultClientConfig() config.Config {
	return config.Config{
		Signing: &config.Signing{
			Default: &config.SigningProfile{
				Expiry:       168 * time.Hour,
				ExpiryString: "168h",
			},
			Profiles: map[string]*config.SigningProfile{
				"client": &config.SigningProfile{
					Expiry:       defaultConfig.Expiry,
					ExpiryString: defaultConfig.ExpiryString,
					Usage:        []string{"signing", "key encipherment", "client auth"},
				},
			},
		},
	}
}
