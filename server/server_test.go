package server

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantMethod string
		wantPath   string
		wantN      int
		wantErr    bool
	}{
		{
			name:       "GET request",
			raw:        "GET example.txt\n",
			wantMethod: "GET",
			wantPath:   "example.txt",
			wantN:      16,
			wantErr:    false,
		},
		{
			name:       "Bad request",
			raw:        "ECHO hello\n",
			wantMethod: "",
			wantPath:   "",
			wantN:      11,
			wantErr:    true,
		},
		{
			name:       "empty request",
			raw:        "\n",
			wantMethod: "",
			wantPath:   "",
			wantN:      1,
			wantErr:    true,
		},
		{
			name:       "empty request no newline",
			raw:        "",
			wantMethod: "",
			wantPath:   "",
			wantN:      0,
			wantErr:    false,
		},
		{
			name:       "Good request no path",
			raw:        "GET \n",
			wantMethod: "GET",
			wantPath:   "",
			wantN:      5,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, n, err := parseRequest([]byte(tt.raw))

			if err != nil && !tt.wantErr {
				t.Fatalf("expected err to be nil, err is %s", err)
			}

			if n != tt.wantN {
				t.Fatalf("expected %v, got %v", tt.wantN, n)
			}

			if got.Method != tt.wantMethod {
				t.Errorf("expected method to be %s, got %s", tt.wantMethod, got.Method)
			}

			if got.Path != tt.wantPath {
				t.Errorf("expected method to be %s, got %s", tt.wantPath, got.Path)
			}
		})
	}
}
