package logz

import (
	"net/http"
	"reflect"
	"testing"
)

func TestMaskName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"John Doe", args{"John Doe"}, "J**n D*e"},
		{"Jo Do", args{"Jo Do"}, "** **"},
		{"John", args{"John"}, "J**n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskName(tt.args.s); got != tt.want {
				t.Errorf("MaskName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskEmail(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test@email.com", args{"test@gmail.com"}, "t**t@gmail.com"},
		{"tt@email.com", args{"tt@gmail.com"}, "**@gmail.com"},
		{"email.com", args{"email.com"}, "email.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskEmail(tt.args.s); got != tt.want {
				t.Errorf("MaskEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskMap(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			"flat map",
			args{map[string]any{"NAME": "John Doe", "email": "john@doe.com"}},
			map[string]any{"NAME": "J**n D*e", "email": "j**n@doe.com"},
		},
		{
			"nested map",
			args{map[string]any{
				"user": map[string]any{
					"name":  "John Doe",
					"EMAIL": "john@doe.com",
				},
			}},
			map[string]any{
				"user": map[string]any{
					"name":  "J**n D*e",
					"EMAIL": "j**n@doe.com",
				},
			},
		},
		{
			"map slice",
			args{map[string]any{
				"users": []any{
					map[string]any{"name": "John Doe"},
				},
			}},
			map[string]any{
				"users": []any{
					map[string]any{"name": "J**n D*e"},
				},
			},
		},
	}

	SetReplacerMap(map[string]func(string) string{
		"name":  MaskName,
		"email": MaskEmail,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaskMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskHttpHeader(t *testing.T) {
	type args struct {
		h http.Header
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			"single value",
			args{http.Header{"SECRET": []string{"secret_value"}}},
			http.Header{"SECRET": []string{"****"}},
		},
		{
			"multi values",
			args{http.Header{"secret": []string{"secret_value_1", "secret_value_2"}}},
			http.Header{"secret": []string{"****"}},
		},
	}

	SetReplacerMap(map[string]func(string) string{
		"secret": Mask,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskHttpHeader(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaskHttpHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
