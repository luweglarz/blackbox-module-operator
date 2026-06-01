/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import "testing"

func TestValidateRejectsInlineSecrets(t *testing.T) {
	tests := []struct {
		name    string
		http    *HTTPProbe
		wantErr bool
	}{
		{
			name:    "no http probe",
			http:    nil,
			wantErr: false,
		},
		{
			name:    "http probe without credentials",
			http:    &HTTPProbe{Method: "GET"},
			wantErr: false,
		},
		{
			name:    "basic auth password file is allowed",
			http:    &HTTPProbe{BasicAuth: &BasicAuth{Username: "u", PasswordFile: "/etc/secret"}},
			wantErr: false,
		},
		{
			name:    "inline basic auth password is rejected",
			http:    &HTTPProbe{BasicAuth: &BasicAuth{Username: "u", Password: "hunter2"}},
			wantErr: true,
		},
		{
			name:    "inline authorization credentials are rejected",
			http:    &HTTPProbe{Authorization: &Authorization{Type: "Bearer", Credentials: "token"}},
			wantErr: true,
		},
		{
			name:    "inline oauth2 client secret is rejected",
			http:    &HTTPProbe{OAuth2: &OAuth2Config{ClientID: "id", ClientSecret: "secret"}},
			wantErr: true,
		},
		{
			name:    "oauth2 client secret ref is allowed",
			http:    &HTTPProbe{OAuth2: &OAuth2Config{ClientID: "id", ClientSecretRef: "my-secret"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := BlackboxModuleSpec{Prober: "http", HTTP: tt.http}
			if tt.http == nil {
				// Keep the prober/config-consistency check satisfied so the test
				// isolates the inline-secret validation.
				spec = BlackboxModuleSpec{Prober: "tcp", TCP: &TCPProbe{}}
			}
			m := &BlackboxModule{Spec: spec}
			_, err := m.validateBlackboxModule()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateBlackboxModule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
