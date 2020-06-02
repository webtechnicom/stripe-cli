package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stripe/stripe-cli/pkg/config"
	"github.com/stripe/stripe-mock/core"
)

type mockCmd struct {
	cmd    *cobra.Command
	config *config.Config

	mockOptions core.Options
}

func newMockCmd() *mockCmd {
	mc := &mockCmd{
		config: &Config,
	}
	mc.cmd = &cobra.Command{
		Use:   "mock",
		Short: "Start a stripe-mock server",
		Long: `This exposes the stripe-mock server as a stripe-cli command.
stripe-mock is a mock HTTP server that responds like the real Stripe API.
It can be used instead of Stripe's test mode to make test suites integrating with Stripe faster and less brittle.`,
		Example: `stripe mock --https-port 1212`,
		RunE:    mc.runMockCmd,
	}

	mc.cmd.Flags().BoolVar(&mc.mockOptions.Http, "http", false, "Run with HTTP")
	mc.cmd.Flags().StringVar(&mc.mockOptions.HttpAddr, "http-addr", "", fmt.Sprintf("Host and port to listen on for HTTP as `<ip>:<port>`; empty <ip> to bind all system IPs, empty <port> to have system choose; e.g. ':%v', '127.0.0.1:%v'", core.DefaultPortHTTP, core.DefaultPortHTTP))
	mc.cmd.Flags().IntVar(&mc.mockOptions.HttpPort, "http-port", -1, "Port to listen on for HTTP; same as '-http-addr :<port>'")
	mc.cmd.Flags().StringVar(&mc.mockOptions.HttpUnixSocket, "http-unix", "", "Unix socket to listen on for HTTP")
	mc.cmd.Flags().BoolVar(&mc.mockOptions.Https, "https", false, "Run with HTTPS; also enables HTTP/2")
	mc.cmd.Flags().StringVar(&mc.mockOptions.HttpsAddr, "https-addr", "", fmt.Sprintf("Host and port to listen on for HTTPS as `<ip>:<port>`; empty <ip> to bind all system IPs, empty <port> to have system choose; e.g. ':%v', '127.0.0.1:%v'", core.DefaultPortHTTPS, core.DefaultPortHTTPS))
	mc.cmd.Flags().IntVar(&mc.mockOptions.HttpsPort, "https-port", -1, "Port to listen on for HTTPS; same as '-https-addr :<port>'")
	mc.cmd.Flags().StringVar(&mc.mockOptions.HttpsUnixSocket, "https-unix", "", "Unix socket to listen on for HTTPS")
	mc.cmd.Flags().IntVar(&mc.mockOptions.Port, "port", -1, "Port to listen on; also respects PORT from environment")
	mc.cmd.Flags().StringVar(&mc.mockOptions.FixturesPath, "fixtures", "", "Path to fixtures to use instead of bundled version (should be JSON)")
	mc.cmd.Flags().StringVar(&mc.mockOptions.SpecPath, "spec", "", "Path to OpenAPI spec to use instead of bundled version (should be JSON)")
	mc.cmd.Flags().BoolVar(&mc.mockOptions.StrictVersionCheck, "strict-version-check", false, "Errors if version sent in Stripe-Version doesn't match the one in OpenAPI")
	mc.cmd.Flags().StringVar(&mc.mockOptions.UnixSocket, "unix", "", "Unix socket to listen on")
	// mc.cmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose mode") // TODO(bwang): what is this for?
	mc.cmd.Flags().BoolVar(&mc.mockOptions.ShowVersion, "version", false, "Show version and exit")

	return mc
}

// This command simply passes its CLI flags through the interface exposed by the stripe-mock library, where the mock server implemention lives.
func (mc *mockCmd) runMockCmd(cmd *cobra.Command, args []string) error {
	// TODO(bwang): There's a minor mismatch in flag syntax stripe-mock and stripe-cli -- so the error messages emitted aren't accurate when calling
	// via the stripe-cli. We could change the error messages in stripe-mock to use double '--' since stripe-mock supports both forms.
	// stripe-mock msgs mention single '-' flags, but stripe-cli uses double '--' flags.
	// ```
	// st-bwang1:stripe-cli bwang$ go run cmd/stripe/main.go mock --port 2 --unix 2
	// Invalid options: Please specify only one of -port or -unix
	// ```

	stripeMock, err := core.NewStripeMock(mc.mockOptions)
	if err != nil {
		return err
	}

	if err := stripeMock.Start(); err != nil {
		return err
	}

	// Block forever. The program will aborted if there are errors with the stripe-mock goroutines serving incoming connections.
	select {}
}
