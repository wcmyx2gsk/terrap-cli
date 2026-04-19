package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistryClient_Defaults(t *testing.T) {
	client := NewRegistryClient("")
	assert.NotNil(t, client)
}

func TestRegistryClient_FetchVersions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"versions": [
				{"version": "3.0.0"},
				{"version": "2.9.1"},
				{"version": "1.0.0"}
			]
		}`))
	}))
	defer server.Close()

	client := NewRegistryClient(server.URL)
	require.NotNil(t, client)

	versions, err := client.FetchVersions("hashicorp", "aws")
	require.NoError(t, err)
	assert.Len(t, versions, 3)
	assert.Contains(t, versions, "3.0.0")
	assert.Contains(t, versions, "2.9.1")
	assert.Contains(t, versions, "1.0.0")
}

func TestRegistryClient_FetchVersions_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewRegistryClient(server.URL)
	require.NotNil(t, client)

	_, err := client.FetchVersions("nonexistent", "provider")
	assert.Error(t, err)
}

func TestRegistryClient_FetchVersions_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	client := NewRegistryClient(server.URL)
	require.NotNil(t, client)

	_, err := client.FetchVersions("hashicorp", "aws")
	assert.Error(t, err)
}

func TestRegistryClient_FetchVersions_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"versions": []}`))
	}))
	defer server.Close()

	client := NewRegistryClient(server.URL)
	require.NotNil(t, client)

	versions, err := client.FetchVersions("hashicorp", "aws")
	require.NoError(t, err)
	assert.Empty(t, versions)
}
