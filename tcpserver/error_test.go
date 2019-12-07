package tcpserver

import (
	"errors"
	"testing"

	"github.com/go-test/deep"
)

func TestError_methods(t *testing.T) {
	err := errors.New("test")
	tests := []struct {
		name    string
		action  string
		err     error
		value   interface{}
		wantErr string
		wantAct string
		wantVal interface{}
	}{
		{"normal", "alRT", err, "value", "test", "alRT", "value"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				action: tt.action,
				err:    tt.err,
				value:  tt.value,
			}
			if got := e.Action(); got != tt.wantAct {
				t.Errorf("Error.Action() = %v, want %v", got, tt.wantAct)
			}

			if got := e.Error(); got != tt.wantErr {
				t.Errorf("Error.Error() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name   string
		action string
		err    error
		want   error
	}{
		{"normal", "alert", errors.New("test"), &Error{action: "alert", err: errors.New("test"), value: nil}},
		{"empty error", "alert", nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.action, tt.err)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewError() %v", diff)
			}
		})
	}
}

func TestNewErrorWithValue(t *testing.T) {
	tests := []struct {
		name   string
		action string
		value  interface{}
		err    error
		want   error
	}{
		{"normal", "alert", 42, errors.New("test"), &Error{action: "alert", err: errors.New("test"), value: 42}},
		{"no value", "alert", nil, errors.New("test"), &Error{action: "alert", err: errors.New("test"), value: nil}},
		{"empty error", "alert", nil, nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrorWithValue(tt.action, tt.value, tt.err)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewError() %v", diff)
			}
		})
	}
}
