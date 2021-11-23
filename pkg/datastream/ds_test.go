package datastream

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func mockAPIClient(t *testing.T, mockServer *httptest.Server) DS {
	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AddCert(mockServer.Certificate())
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	s, err := session.New(
		session.WithClient(httpClient),
		session.WithSigner(&edgegrid.Config{Host: serverURL.Host}),
	)
	assert.NoError(t, err)
	return Client(s)
}

func TestClient(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *ds
	}{
		"no options provided, return default": {
			options: nil,
			expected: &ds{
				Session: sess,
			},
		},
		"option provided, overwrite session": {
			options: []Option{func(c *ds) {
				c.Session = nil
			}},
			expected: &ds{
				Session: nil,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess, test.options...)
			assert.Equal(t, res, test.expected)
		})
	}
}
