package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func sendTestRequest(t *testing.T, method, url string) (int, string) {
	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
