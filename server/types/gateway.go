package types

import (
	"net/url"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/flashbots/go-boost-utils/utils"
)

// Gateway represents a relay that mev-boost connects to.
type Gateway struct {
	PublicKey phase0.BLSPubKey
	URL       *url.URL
}

func (g *Gateway) String() string {
	return g.URL.String()
}

// GetURI returns the full request URI with scheme, host, path and args for the gateway.
func (g *Gateway) GetURI(path string) string {
	return GetURI(g.URL, path)
}

// NewGateway creates a new instance based on an input string
// relayURL can be IP@PORT, PUBKEY@IP:PORT, https://IP, etc.
func NewGateway(relayURL string) (gateway Gateway, err error) {
	// Add protocol scheme prefix if it does not exist.
	if !strings.HasPrefix(relayURL, "http") {
		relayURL = "http://" + relayURL
	}

	// Parse the provided relay's URL and save the parsed URL in the Gateway.
	gateway.URL, err = url.ParseRequestURI(relayURL)
	if err != nil {
		return gateway, err
	}

	// Extract the relay's public key from the parsed URL.
	if gateway.URL.User.Username() == "" {
		return gateway, ErrMissingRelayPubkey
	}

	// Convert the username string to a public key.
	gateway.PublicKey, err = utils.HexToPubkey(gateway.URL.User.Username())
	if err != nil {
		return gateway, err
	}

	// Check if the public key is the point-at-infinity.
	if gateway.PublicKey.IsInfinity() {
		return gateway, ErrPointAtInfinityPubkey
	}

	return gateway, nil
}

// GatewayToString returns the string representation of a gateway entry
func GatewayToString(gateway Gateway) string {
	return gateway.String()
}
