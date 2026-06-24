package options

import (
	"time"

	"github.com/spf13/pflag"
)

const defaultOAuthStateTTL = 5 * time.Minute

// OAuthOptions defines options for the generic OAuth login flow.
type OAuthOptions struct {
	AuthorizationEndpoint string        `json:"authorization-endpoint" mapstructure:"authorization-endpoint"`
	ClientID              string        `json:"client-id"              mapstructure:"client-id"`
	RedirectURI           string        `json:"redirect-uri"           mapstructure:"redirect-uri"`
	Scopes                []string      `json:"scopes"                 mapstructure:"scopes"`
	StateTTL              time.Duration `json:"state-ttl"              mapstructure:"state-ttl"`
	CookieSecure          bool          `json:"cookie-secure"          mapstructure:"cookie-secure"`
}

// NewOAuthOptions creates OAuth options with safe defaults.
func NewOAuthOptions() *OAuthOptions {
	return &OAuthOptions{
		Scopes:       []string{"openid", "profile", "email"},
		StateTTL:     defaultOAuthStateTTL,
		CookieSecure: true,
	}
}

// Validate verifies OAuth options.
func (o *OAuthOptions) Validate() []error {
	return nil
}

// AddFlags adds OAuth-related command-line flags.
func (o *OAuthOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&o.AuthorizationEndpoint,
		"oauth.authorization-endpoint",
		o.AuthorizationEndpoint,
		"Keycloak OpenID Connect authorization endpoint.",
	)
	fs.StringVar(&o.ClientID, "oauth.client-id", o.ClientID, "OAuth client ID.")
	fs.StringVar(&o.RedirectURI, "oauth.redirect-uri", o.RedirectURI, "OAuth callback URI.")
	fs.StringSliceVar(&o.Scopes, "oauth.scopes", o.Scopes, "OAuth scopes requested from the identity provider.")
	fs.DurationVar(&o.StateTTL, "oauth.state-ttl", o.StateTTL, "Maximum lifetime of the OAuth state cookie.")
	fs.BoolVar(
		&o.CookieSecure,
		"oauth.cookie-secure",
		o.CookieSecure,
		"Only send the OAuth state cookie over HTTPS.",
	)
}
