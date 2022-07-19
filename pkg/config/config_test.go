package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		path        string
		want        Config
		shouldThrow bool
	}{
		{path: ".nonexistentenv", shouldThrow: true},
		{path: "../../.example.env", want: Config{Login: "foo@bar.com", Token: "token", Port: "8080"}, shouldThrow: false},
	}

	for _, tc := range tests {
		if err := Read(tc.path); err != nil && !tc.shouldThrow {
			t.Errorf("Unexpected error while reading config %s", tc.path)
		}
		cfg := Get()
		if !cmp.Equal(&tc.want, cfg) {
			t.Logf("%+v", cfg)
			t.Errorf("Read wrong values from config %s", tc.path)
		}
	}
}
