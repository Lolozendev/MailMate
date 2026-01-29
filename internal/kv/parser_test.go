package kv

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "simple key-value pairs",
			input: "Name='John Doe';Count=5",
			want: map[string]string{
				"Name":  "John Doe",
				"Count": "5",
			},
			wantErr: false,
		},
		{
			name:  "double quotes",
			input: `Name="Jane Doe";Email="jane@example.com"`,
			want: map[string]string{
				"Name":  "Jane Doe",
				"Email": "jane@example.com",
			},
			wantErr: false,
		},
		{
			name:  "mixed quotes and unquoted",
			input: "Name='John';Count=42;Active=true",
			want: map[string]string{
				"Name":   "John",
				"Count":  "42",
				"Active": "true",
			},
			wantErr: false,
		},
		{
			name:  "whitespace handling",
			input: " Name = 'John' ; Count = 5 ",
			want: map[string]string{
				"Name":  "John",
				"Count": "5",
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "missing equals sign",
			input:   "Name'John'",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   "='value'",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnquote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single quotes",
			input: "'hello'",
			want:  "hello",
		},
		{
			name:  "double quotes",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "no quotes",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "mismatched quotes",
			input: `'hello"`,
			want:  `'hello"`,
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unquote(tt.input); got != tt.want {
				t.Errorf("unquote() = %v, want %v", got, tt.want)
			}
		})
	}
}
