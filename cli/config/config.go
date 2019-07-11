package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddr         string        `envconfig:"SERVER_ADDR" default:"localhost:8888"`
	RequestFile        string        `envconfig:"REQUEST_FILE"`
	PrintSampleRequest bool          `envconfig:"PRINT_SAMPLE_REQUEST"`
	ResponseFormat     string        `envconfig:"RESPONSE_FORMAT" default:"json"`
	Timeout            time.Duration `envconfig:"TIMEOUT" default:"10s"`
	TLS                bool          `envconfig:"TLS"`
	ServerName         string        `envconfig:"TLS_SERVER_NAME"`
	InsecureSkipVerify bool          `envconfig:"TLS_INSECURE_SKIP_VERIFY"`
	CACertFile         string        `envconfig:"TLS_CA_CERT_FILE"`
	CertFile           string        `envconfig:"TLS_CERT_FILE"`
	KeyFile            string        `envconfig:"TLS_KEY_FILE"`
	AuthToken          string        `envconfig:"AUTH_TOKEN"`
	AuthTokenType      string        `envconfig:"AUTH_TOKEN_TYPE" default:"Bearer"`
	JWTKey             string        `envconfig:"JWT_KEY"`
	JWTKeyFile         string        `envconfig:"JWT_KEY_FILE"`
}

var C = NewConfig()

func NewConfig() *Config {
	c := &Config{}
	envconfig.Process("", c)
	return c
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.ServerAddr, "server-addr", "s", c.ServerAddr, "server address in form of host:port")
	fs.StringVarP(&c.RequestFile, "request-file", "f", c.RequestFile, "client request file (must be json, yaml, or xml); use \"-\" for stdin + json")
	fs.BoolVarP(&c.PrintSampleRequest, "print-sample-request", "p", c.PrintSampleRequest, "print sample request file and exit")
	fs.StringVarP(&c.ResponseFormat, "response-format", "o", c.ResponseFormat, "response format (json, prettyjson, yaml, or xml)")
	fs.DurationVar(&c.Timeout, "timeout", c.Timeout, "client connection timeout")
	fs.BoolVar(&c.TLS, "tls", c.TLS, "enable tls")
	fs.StringVar(&c.ServerName, "tls-server-name", c.ServerName, "tls server name override")
	fs.BoolVar(&c.InsecureSkipVerify, "tls-insecure-skip-verify", c.InsecureSkipVerify, "INSECURE: skip tls checks")
	fs.StringVar(&c.CACertFile, "tls-ca-cert-file", c.CACertFile, "ca certificate file")
	fs.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, "client certificate file")
	fs.StringVar(&c.KeyFile, "tls-key-file", c.KeyFile, "client key file")
	fs.StringVar(&c.AuthToken, "auth-token", c.AuthToken, "authorization token")
	fs.StringVar(&c.AuthTokenType, "auth-token-type", c.AuthTokenType, "authorization token type")
	fs.StringVar(&c.JWTKey, "jwt-key", c.JWTKey, "jwt key")
	fs.StringVar(&c.JWTKeyFile, "jwt-key-file", c.JWTKeyFile, "jwt key file")
}
