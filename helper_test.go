package tiktok_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/jiutian1000/tiktok"
	"github.com/stretchr/testify/require"
)

type TestRecord struct {
	Name    string          `json:"name"`
	Args    json.RawMessage `json:"args"`
	Request struct {
		Method string          `json:"method"`
		URL    string          `json:"url"`
		Query  string          `json:"query"`
		Body   json.RawMessage `json:"body"`
	} `json:"request"`
	Response struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    json.RawMessage   `json:"body"`
	} `json:"response"`
	Want    json.RawMessage `json:"want"`
	WantErr bool            `json:"want_err"`
}

func loadTestData(t *testing.T, fn string) (records []TestRecord) {
	t.Helper()
	b, err := ioutil.ReadFile(fn)
	require.NoError(t, err)
	err = json.Unmarshal(b, &records)
	require.NoError(t, err)
	return
}

func setupMock(t *testing.T, tt TestRecord, args, want interface{}) {
	httpmock.RegisterResponder(
		tt.Request.Method, tt.Request.URL,
		func(r *http.Request) (res *http.Response, err error) {
			if tt.Request.Query != "" {
				require.Equal(t, tt.Request.Query, r.URL.RawQuery)
			}

			if len(tt.Request.Body) > 0 {
				defer r.Body.Close()
				b, _ := io.ReadAll(r.Body)
				require.JSONEq(t, string(tt.Request.Body), string(b))
			} else {
				tt.Request.Body = []byte(`{}`)
			}
			res, err = httpmock.NewJsonResponse(tt.Response.Status, tt.Response.Body)
			return
		},
	)

	var err error
	if want != nil {
		err = json.Unmarshal([]byte(tt.Want), want)
		require.NoError(t, err)
	}
	if args != nil {
		err = json.Unmarshal(tt.Args, args)
		require.NoError(t, err)
	}
}

func mockTime() func() {
	tiktok.Timestamp = func() string {
		return "1600000000"
	}

	return func() {
		tiktok.Timestamp = func() string {
			return fmt.Sprintf("%d", time.Now().Unix())
		}
	}
}

func mockTests(t *testing.T, fn string, args, want interface{}, exec func() (interface{}, error)) {
	t.Helper()

	restore := mockTime()
	defer restore()

	tests := loadTestData(t, fn)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			setupMock(t, tt, args, want)

			got, err := exec()
			if tt.WantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				if want != nil {
					require.Equal(t, want, got)
				}
			}
		})
	}
}
