package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/kmova/waytogo/api"
	cliflags "github.com/kmova/waytogo/cli/flags"
	"github.com/kmova/waytogo/cliconfig"
	"github.com/kmova/waytogo/cliconfig/configfile"
	"github.com/kmova/waytogo/version"
	"github.com/kmova/waytogo/opts"
	"github.com/kmova/waytogo/pkg/term"
	"github.com/kmova/waytogo/client"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
)

// WaytogoCli represents the waytogo command line client.
// Instances of the client can be returned from NewWaytogoCli.
type WaytogoCli struct {
	// initializing closure
	init func() error

	// configFile has the client configuration file
	configFile *configfile.ConfigFile
	// in holds the input stream and closer (io.ReadCloser) for the client.
	in io.ReadCloser
	// out holds the output stream (io.Writer) for the client.
	out io.Writer
	// err holds the error stream (io.Writer) for the client.
	err io.Writer
	// keyFile holds the key file as a string.
	keyFile string
	// inFd holds the file descriptor of the client's STDIN (if valid).
	inFd uintptr
	// outFd holds file descriptor of the client's STDOUT (if valid).
	outFd uintptr
	// isTerminalIn indicates whether the client's STDIN is a TTY
	isTerminalIn bool
	// isTerminalOut indicates whether the client's STDOUT is a TTY
	isTerminalOut bool
	// client is the http client that performs all API operations
	client client.APIClient
	// state holds the terminal input state
	inState *term.State
	// outState holds the terminal output state
	outState *term.State
}

// Initialize calls the init function that will setup the configuration for the client
// such as the TLS, tcp and other parameters used to run the client.
func (cli *WaytogoCli) Initialize() error {
	if cli.init == nil {
		return nil
	}
	return cli.init()
}

// Client returns the APIClient
func (cli *WaytogoCli) Client() client.APIClient {
	return cli.client
}

// Out returns the writer used for stdout
func (cli *WaytogoCli) Out() io.Writer {
	return cli.out
}

// Err returns the writer used for stderr
func (cli *WaytogoCli) Err() io.Writer {
	return cli.err
}

// In returns the reader used for stdin
func (cli *WaytogoCli) In() io.ReadCloser {
	return cli.in
}

// ConfigFile returns the ConfigFile
func (cli *WaytogoCli) ConfigFile() *configfile.ConfigFile {
	return cli.configFile
}

// IsTerminalOut returns true if the clients stdin is a TTY
func (cli *WaytogoCli) IsTerminalOut() bool {
	return cli.isTerminalOut
}

// OutFd returns the fd for the stdout stream
func (cli *WaytogoCli) OutFd() uintptr {
	return cli.outFd
}

// CheckTtyInput checks if we are trying to attach to a container tty
// from a non-tty client input stream, and if so, returns an error.
func (cli *WaytogoCli) CheckTtyInput(attachStdin, ttyMode bool) error {
	// In order to attach to a container tty, input stream for the client must
	// be a tty itself: redirecting or piping the client standard input is
	// incompatible with `waytogo run -t`, `waytogo exec -t` or `waytogo attach`.
	if ttyMode && attachStdin && !cli.isTerminalIn {
		eText := "the input device is not a TTY"
		if runtime.GOOS == "windows" {
			return errors.New(eText + ".  If you are using mintty, try prefixing the command with 'winpty'")
		}
		return errors.New(eText)
	}
	return nil
}

// PsFormat returns the format string specified in the configuration.
// String contains columns and format specification, for example {{ID}}\t{{Name}}.
func (cli *WaytogoCli) PsFormat() string {
	return cli.configFile.PsFormat
}

// ImagesFormat returns the format string specified in the configuration.
// String contains columns and format specification, for example {{ID}}\t{{Name}}.
func (cli *WaytogoCli) ImagesFormat() string {
	return cli.configFile.ImagesFormat
}

func (cli *WaytogoCli) setRawTerminal() error {
	if os.Getenv("NORAW") == "" {
		if cli.isTerminalIn {
			state, err := term.SetRawTerminal(cli.inFd)
			if err != nil {
				return err
			}
			cli.inState = state
		}
		if cli.isTerminalOut {
			state, err := term.SetRawTerminalOutput(cli.outFd)
			if err != nil {
				return err
			}
			cli.outState = state
		}
	}
	return nil
}

func (cli *WaytogoCli) restoreTerminal(in io.Closer) error {
	if cli.inState != nil {
		term.RestoreTerminal(cli.inFd, cli.inState)
	}
	if cli.outState != nil {
		term.RestoreTerminal(cli.outFd, cli.outState)
	}
	// WARNING: DO NOT REMOVE THE OS CHECK !!!
	// For some reason this Close call blocks on darwin..
	// As the client exists right after, simply discard the close
	// until we find a better solution.
	if in != nil && runtime.GOOS != "darwin" {
		return in.Close()
	}
	return nil
}

// NewWaytogoCli returns a WaytogoCli instance with IO output and error streams set by in, out and err.
// The key file, protocol (i.e. unix) and address are passed in as strings, along with the tls.Config. If the tls.Config
// is set the client scheme will be set to https.
// The client will be given a 32-second timeout (see https://github.com/docker/docker/pull/8035).
func NewWaytogoCli(in io.ReadCloser, out, err io.Writer, clientFlags *cliflags.ClientFlags) *WaytogoCli {
	cli := &WaytogoCli{
		in:      in,
		out:     out,
		err:     err,
	}

	cli.init = func() error {
		clientFlags.PostParse()
		cli.configFile = LoadDefaultConfigFile(err)

		client, err := NewAPIClientFromFlags(clientFlags, cli.configFile)
		if err != nil {
			return err
		}

		cli.client = client

		if cli.in != nil {
			cli.inFd, cli.isTerminalIn = term.GetFdInfo(cli.in)
		}
		if cli.out != nil {
			cli.outFd, cli.isTerminalOut = term.GetFdInfo(cli.out)
		}

		return nil
	}

	return cli
}

// LoadDefaultConfigFile attempts to load the default config file and returns
// an initialized ConfigFile struct if none is found.
func LoadDefaultConfigFile(err io.Writer) *configfile.ConfigFile {
	configFile, e := cliconfig.Load(cliconfig.ConfigDir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	return configFile
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
func NewAPIClientFromFlags(clientFlags *cliflags.ClientFlags, configFile *configfile.ConfigFile) (client.APIClient, error) {
	host, err := getServerHost(clientFlags.Common.Hosts, nil) //TODO : kiran
	if err != nil {
		return &client.Client{}, err
	}

	customHeaders := configFile.HTTPHeaders
	if customHeaders == nil {
		customHeaders = map[string]string{}
	}
	customHeaders["User-Agent"] = clientUserAgent()

	verStr := api.DefaultVersion
	if tmpStr := os.Getenv("WAYTOGO_API_VERSION"); tmpStr != "" {
		verStr = tmpStr
	}

	httpClient, err := newHTTPClient(host, nil ) //TODO: kiran
	if err != nil {
		return &client.Client{}, err
	}

	return client.NewClient(host, verStr, httpClient, customHeaders)
}

func getServerHost(hosts []string, tlsOptions *tlsconfig.Options) (host string, err error) {
	switch len(hosts) {
	case 0:
		host = os.Getenv("WAYTOGO_HOST")
	case 1:
		host = hosts[0]
	default:
		return "", errors.New("Please specify only one -H")
	}

	host, err = opts.ParseHost(tlsOptions != nil, host)
	return
}

func newHTTPClient(host string, tlsOptions *tlsconfig.Options) (*http.Client, error) {
	if tlsOptions == nil {
		// let the api client configure the default transport.
		return nil, nil
	}

	config, err := tlsconfig.Client(*tlsOptions)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	proto, addr, _, err := client.ParseHost(host)
	if err != nil {
		return nil, err
	}

	sockets.ConfigureTransport(tr, proto, addr)

	return &http.Client{
		Transport: tr,
	}, nil
}

func clientUserAgent() string {
	return "WayToGo-Client/" + version.Version + " (" + runtime.GOOS + ")"
}
