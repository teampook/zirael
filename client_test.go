package zirael

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetSuccess(t *testing.T)  {
	client := NewClient("", "", "", WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)

	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")
	response, err := client.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestPostSuccess(t *testing.T)  {
	client := NewClient("", "", "", WithHTTPTimeout(10 * time.Millisecond))

	requestBodyString := `{ "hello": "Bareksa" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Post(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer response.Body.Close()
	}

	respBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")

	return string(respBody)
}