package config

type DrandConfig struct {
	Servers      []string
	GossipRelays []string
	ChainInfo    string
}

func NewDrandConfig() *DrandConfig {
	cfg := &DrandConfig{
		Servers: []string{
			"https://api.drand.sh",
			"https://api2.drand.sh",
			"https://api3.drand.sh",
			"https://drand.cloudflare.com",
		},
		GossipRelays: []string{
			"/dnsaddr/api.drand.sh/",
			"/dnsaddr/api2.drand.sh/",
			"/dnsaddr/api3.drand.sh/",
		},
		ChainInfo: `{
			"public_key": "868f005eb8e6e4ca0a47c8a77ceaa5309a47978a7c71bc5cce96366b5d7a569937c529eeda66c7293784a9402801af31",
			"period": 30,
			"genesis_time": 1595431050,
			"hash": "8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce"
		}`,
	}
	return cfg
}
