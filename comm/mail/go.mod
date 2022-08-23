module mail

go 1.18

replace github.com/oldbai555/log => ../../log

replace github.com/oldbai555/comm => ../../comm

require (
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-imap-id v0.0.0-20190926060100-f94a56b9ecde
	github.com/emersion/go-message v0.16.0
	github.com/oldbai555/comm v0.0.0-20220806102209-f378817e4c79
	github.com/oldbai555/log v0.0.0-00010101000000-000000000000
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/elliotchance/pie v1.39.0 // indirect
	github.com/emersion/go-sasl v0.0.0-20200509203442-7bfe0ed36a21 // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/forgoer/openssl v1.2.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)
