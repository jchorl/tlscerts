module github.com/jchorl/tlscerts

go 1.13

require (
	github.com/cloudflare/cfssl v1.4.1
	honnef.co/go/tools v0.0.1-2019.2.3
)

replace github.com/cloudflare/cfssl => github.com/jchorl/cfssl v1.4.2-0.20191120163807-27b89981a159
