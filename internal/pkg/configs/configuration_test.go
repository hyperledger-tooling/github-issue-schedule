package configs

import (
	"testing"
)

func TestConfigurationValidate(t *testing.T) {
	cases := []struct {
		desc   string
		conf   Configuration
		broken bool
	}{
		{"negative", Configuration{BufferWindowDays: -1}, true},
		{"too large", Configuration{BufferWindowDays: 100001}, true},
		{"zero", Configuration{BufferWindowDays: 0}, false},
		{"normal", Configuration{BufferWindowDays: 7}, false},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			err := c.conf.validate()
			if c.broken && err == nil {
				t.Fatalf("expected error for config %+v", c.conf)
			}
			if !c.broken && err != nil {
				t.Fatalf("did not expect error: %v", err)
			}
		})
	}
}
