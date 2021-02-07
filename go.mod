module github.com/jchorl/tlscerts

go 1.15

require (
	github.com/cloudflare/cfssl v1.4.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/cloudflare/cfssl => github.com/jchorl/cfssl v1.4.2-0.20191120163807-27b89981a159
