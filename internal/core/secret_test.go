package core_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Untanky/modern-auth/internal/core"
)

type demoJson struct {
	Secret *core.SecretValue `json:"some_value"`
}

func (d *demoJson) String() string {
	return d.Secret.String()
}

func TestSecretMarshalling(t *testing.T) {
	tests := []struct {
		name     string
		secret   fmt.Stringer
		wantJSON string
	}{
		{
			name:     "Test Secret JSON Marshalling",
			secret:   core.NewSecretValue("secret"),
			wantJSON: "\"y25ORpgIz3M3ZAB2BQywdR1CA3b/akE+js4KPz8EqB7K+AnAa1qY8DKXrdv65NlAfwNEGhaDKvItU1Vp27tAkg==\"",
		},
		{
			name:     "Test Secret JSON Marshalling in struct",
			secret:   &demoJson{Secret: core.NewSecretValue("secret")},
			wantJSON: "{\"some_value\":\"y25ORpgIz3M3ZAB2BQywdR1CA3b/akE+js4KPz8EqB7K+AnAa1qY8DKXrdv65NlAfwNEGhaDKvItU1Vp27tAkg==\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretJSON, err := json.Marshal(tt.secret)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(secretJSON) != tt.wantJSON {
				t.Errorf("MarshalJSON() got = %v, want %v", string(secretJSON), tt.wantJSON)
			}

			err = json.Unmarshal(secretJSON, tt.secret)
			if err == nil {
				t.Errorf("UnmarshalJSON() expected error, but got nil")
				return
			}
		})
	}
}
