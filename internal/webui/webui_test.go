package webui

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	e := store.NewEngine()
	t.Cleanup(e.Stop)
	r := commands.NewRegistry()
	s := New(":0", e, r)
	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func doCommand(t *testing.T, ts *httptest.Server, cmd string) commandResponse {
	t.Helper()
	body, err := json.Marshal(commandRequest{Command: cmd})
	require.NoError(t, err)
	res, err := http.Post(ts.URL+"/api/command", "application/json", strings.NewReader(string(body)))
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var out commandResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&out))
	return out
}

func TestCommandSet(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, "SET foo bar")
	assert.Equal(t, "simple", out.Result["kind"])
	assert.Equal(t, "OK", out.Result["value"])
	assert.GreaterOrEqual(t, out.DurationMs, 0.0)
}

func TestCommandGet(t *testing.T) {
	ts := newTestServer(t)
	doCommand(t, ts, "SET foo bar")
	out := doCommand(t, ts, "GET foo")
	assert.Equal(t, "bulk", out.Result["kind"])
	assert.Equal(t, "bar", out.Result["value"])
}

func TestCommandLowercase(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, "set foo bar")
	assert.Equal(t, "simple", out.Result["kind"])
	got := doCommand(t, ts, "get foo")
	assert.Equal(t, "bar", got.Result["value"])
}

func TestCommandQuotedArgs(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, `SET greeting "hello world"`)
	assert.Equal(t, "OK", out.Result["value"])
	got := doCommand(t, ts, "GET greeting")
	assert.Equal(t, "hello world", got.Result["value"])
}

func TestCommandUnknown(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, "NOSUCHCMD x y")
	assert.Equal(t, "error", out.Result["kind"])
	assert.Contains(t, out.Result["value"], "unknown command")
}

func TestCommandEmpty(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, "   ")
	assert.Equal(t, "error", out.Result["kind"])
}

func TestKeysList(t *testing.T) {
	ts := newTestServer(t)
	doCommand(t, ts, "SET name alice")
	doCommand(t, ts, "RPUSH tasks a b c")
	doCommand(t, ts, "SADD tags go redis")

	res, err := http.Get(ts.URL + "/api/keys")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var out keysResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&out))
	assert.Equal(t, 3, out.Count)

	byKey := map[string]keyInfo{}
	for _, k := range out.Keys {
		byKey[k.Key] = k
	}
	assert.Equal(t, "string", byKey["name"].Type)
	assert.Equal(t, "list", byKey["tasks"].Type)
	assert.Equal(t, "set", byKey["tags"].Type)
}

func TestKeysExpired(t *testing.T) {
	ts := newTestServer(t)
	doCommand(t, ts, "SET temp gone EX 1")

	// Wait for the key to expire plus active-expiration slack.
	time.Sleep(1600 * time.Millisecond)

	res, err := http.Get(ts.URL + "/api/keys")
	require.NoError(t, err)
	defer res.Body.Close()
	var out keysResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&out))
	assert.Equal(t, 0, out.Count)
}

func TestIndexServed(t *testing.T) {
	ts := newTestServer(t)
	res, err := http.Get(ts.URL + "/")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, string(data), "kivo")
}

func TestUnknownPathFallsBackToIndex(t *testing.T) {
	ts := newTestServer(t)
	res, err := http.Get(ts.URL + "/some/deep/link")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, string(data), "<div id=\"root\">")
}

func TestSplitCommand(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"SET foo bar", []string{"SET", "foo", "bar"}},
		{`SET foo "hello world"`, []string{"SET", "foo", "hello world"}},
		{`SET "a b" "c d"`, []string{"SET", "a b", "c d"}},
		{"  PING  ", []string{"PING"}},
		{"", nil},
		{"SADD tag \"with spaces\" plain", []string{"SADD", "tag", "with spaces", "plain"}},
	}
	for _, tc := range cases {
		got := splitCommand(tc.in)
		assert.Equal(t, tc.want, got, "input: %q", tc.in)
	}
}

func TestGetMissingKeyReturnsNull(t *testing.T) {
	ts := newTestServer(t)
	out := doCommand(t, ts, "GET missing")
	assert.Equal(t, "bulk", out.Result["kind"])
	assert.Nil(t, out.Result["value"])
}

func TestCommandInvalidMethod(t *testing.T) {
	ts := newTestServer(t)
	res, err := http.Get(ts.URL + "/api/command")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
