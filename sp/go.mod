module github.com/mitchell/selfpass/sp

go 1.12

require (
	github.com/atotto/clipboard v0.1.2
	github.com/c-bata/go-prompt v0.2.3
	github.com/google/uuid v1.1.1
	github.com/mattn/go-tty v0.0.0-20190424173100-523744f04859 // indirect
	github.com/mitchell/selfpass/services v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/ncw/rclone v1.48.0
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/pquerna/otp v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	gopkg.in/AlecAivazis/survey.v1 v1.8.5
)

replace github.com/mitchell/selfpass/services => ../services

replace github.com/mitchell/selfpass/protobuf/go => ../protobuf/go
