module github.com/oldbai555/comm

go 1.18

replace github.com/oldbai555/comm/extrpkg => ./extrpkg

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/forgoer/openssl v1.2.1
	github.com/oldbai555/comm/extrpkg v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
)

require (
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)
