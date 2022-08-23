module github.com/oldbai555/example

go 1.18

replace (
	github.com/oldbai555/comm => ../comm
	github.com/oldbai555/comm/extrpkg => ../comm/extrpkg
	github.com/oldbai555/comm/mail => ../comm/mail
	github.com/oldbai555/log => ../log
	github.com/oldbai555/web => ../web
)

require (
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-message v0.16.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.11.0
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/oldbai555/comm v0.0.0-20220806102209-f378817e4c79
	github.com/oldbai555/comm/extrpkg v0.0.0-00010101000000-000000000000
	github.com/oldbai555/comm/mail v0.0.0-00010101000000-000000000000
	github.com/oldbai555/log v0.0.0-00010101000000-000000000000
	github.com/oldbai555/web v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.12.0
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/elliotchance/pie v1.39.0 // indirect
	github.com/emersion/go-imap-id v0.0.0-20190926060100-f94a56b9ecde // indirect
	github.com/emersion/go-sasl v0.0.0-20200509203442-7bfe0ed36a21 // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/forgoer/openssl v1.2.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)
