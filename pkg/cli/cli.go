package cli

import "fmt"

// Args represents the command-line arguments for sqlweb.
type Args struct {
	Port       int
	Log        bool
	Help       string
	Version    string
	Connection string
}

// NewArgs initializes and returns a new Args struct with default values.
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

// ValidatePortRange checks if the Port field falls within a valid port number range.
// It returns an error if the port number is invalid.
func (args *Args) ValidatePortRange() error {
	if args.Port < 1 || args.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}
