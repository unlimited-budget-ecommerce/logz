package logz

import (
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
