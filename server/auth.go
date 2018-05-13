package server

import (
	cryptorand "crypto/rand"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

// Parses the token from the header.  If the token could not be parsed,
// returns ("", false), logs an errors, and writes to the response writer.
// Otherwise returns (token, true) with no side-effects.
func parseAuthenticationToken(ctx context.Context, h http.Header, w http.ResponseWriter) (string, bool) {
	header := h.Get("Authentication")
	if header == "" {
		log.Errorf(ctx, "missing Authentication header")
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	if parts[0] != "Token" {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		log.Errorf(ctx, "header %q has wrong format", header)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	return tok, true
}

func generateAPIToken() (string, error) {
	b := make([]byte, 16)
	_, err := cryptorand.Read(b)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read rand")
	}
	return string(hexencode(b)), nil
}
