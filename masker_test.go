package logz

import (
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
			args{map[string]any{"name": "John Doe", "email": "john@doe.com"}},
			map[string]any{"name": "J**n D*e", "email": "j**n@doe.com"},
		},
		{
			"nested map",
			args{map[string]any{
				"user": map[string]any{
					"name":  "John Doe",
					"email": "john@doe.com",
				},
			}},
			map[string]any{
				"user": map[string]any{
					"name":  "J**n D*e",
					"email": "j**n@doe.com",
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
