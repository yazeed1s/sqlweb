package cli

import "fmt"

type Args struct {
	Port       int
	Log        bool
	Help       string
	Version    string
	Connection string
}

func NewArgs() *Args {
	return &Args{
		Port: 3000,
		Log:  false,
		Help: `
			Help information:
			USAGE: sqlweb [OPTION]
			OPTION:
			  -p <port>   	Set the port number (default: 3000)
			  -l=<bool>   	Enable logging (default: false)
			  -h          	Display help information
			  -v          	Display version
			  -c=<schema> 	Use saved connection 
			`,
		Version:    "version 0.1.0",
		Connection: "",
	}
}

func (args *Args) ValidatePortRange() error {
	if args.Port < 1 || args.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}
