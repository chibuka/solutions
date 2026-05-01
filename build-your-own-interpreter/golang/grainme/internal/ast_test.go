package internal

import "testing"

func TestExprString(t *testing.T) {
	tests := []struct {
		e        Expr
		expected string
	}{
		{Literal{value: true}, "true"},
		{Literal{value: false}, "false"},
		{Literal{value: 18}, "18.0"},
		{Literal{value: "bar"}, "bar"},
		{Literal{value: nil}, "nil"},
		{Grouping{inner: Literal{value: "foo"}}, "(group foo)"},
		{Grouping{inner: Grouping{inner: Literal{1}}}, "(group (group 1.0))"},
	}

	for _, tt := range tests {
		got := tt.e.String()
		if got != tt.expected {
			t.Errorf("expected %s got %s", tt.expected, got)
		}
	}
}
