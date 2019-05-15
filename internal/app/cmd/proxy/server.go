package proxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"net"
)

// Server holds the config and CLI references.
type Server struct {
	Metrics *metrics.Metrics
	Config  *config.Config
	CLI     ui.CLI
}

// NewServer holds a configuration and command line interface reference (for log out, etc.)
func NewServer(config *config.Config, metrics *metrics.Metrics) (s Server) {
	s.Config = config
	s.CLI = ui.NewCLI(config)
	s.Metrics = metrics
	return
}

// Server with no params will show the help
func (s *Server) ServerHelp(cmd *cobra.Command, args []string) {
	cli := ui.NewCLI(s.Config)
	fmt.Println(cli.Render("serverUsage", nil))
	return
}

// Start will configure a server and start it.
func (s *Server) Start(cmd *cobra.Command, args []string) {
	Start(s.Config, s.Metrics)
	return

}

// Stop will signal the server to stop.
func (s *Server) Stop(cmd *cobra.Command, args []string) {
	Stop(s.Config, s.Metrics)
	return
}

func PortAvailable(port string) bool {
	host := ":" + port
	server, err := net.Listen("tcp", host)
	if err != nil {
		return false
	}
	server.Close()
	return true
}


