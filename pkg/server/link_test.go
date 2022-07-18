package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// This test doesn't cover any real bitly scenarios
// It just checks that data is processed and parsed correctly given working API
func TestShortenLink(t *testing.T) {
	tests := []struct {
		link        string
		want        string
		shouldThrow bool
	}{
		{
			link:        "https://google.com",
			want:        "https://bit.ly/foo",
			shouldThrow: false,
		},
	}
	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			form, _ := json.Marshal(map[string]string{"link": tc.want})
			rw.Write(form)
		}))
		defer server.Close()

		testBitlyAPI := API{Client: server.Client(), BaseURL: server.URL}
		shortenedLink, err := testBitlyAPI.shortenLink("https")
		if err != nil && !tc.shouldThrow {
			t.Errorf("Unexpected error at %+v", tc)
		}
		if shortenedLink != tc.want {
			t.Errorf("Wanted %s. Got %s", tc.want, shortenedLink)
		}
	}
}
