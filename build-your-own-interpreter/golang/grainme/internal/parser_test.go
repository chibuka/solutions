package internal

import "testing"

func TestParser(t *testing.T) {
	tests := []struct {
		tokens []Token
		want   string
	}{
		{
			tokens: []Token{
				Token{TRUE, "true", ""},
			},
			want: "true",
		},
		{
			tokens: []Token{
				Token{STRING, "bar", "bar"},
			},
			want: "bar",
		},
		{
			tokens: []Token{
				Token{LEFT_PAREN, "(", ""},
				Token{STRING, "foo", "foo"},
				Token{RIGHT_PAREN, ")", ""},
			},
			want: "(group foo)",
		},
		{
			tokens: []Token{
				Token{LEFT_PAREN, "(", ""},
				Token{LEFT_PAREN, "(", ""},
				Token{NUMBER, "42", "42.0"},
				Token{RIGHT_PAREN, ")", ""},
				Token{RIGHT_PAREN, ")", ""},
			},
			want: "(group (group 42.0))",
		},
		{
			tokens: []Token{
				Token{BANG, "!", ""},
				Token{TRUE, "true", ""},
			},
			want: "(! true)",
		},
	}

	for _, tt := range tests {
		parser := NewParse(tt.tokens)
		got, err := parser.Parse()
		if err != nil {
			t.Errorf("got an error: %v", err)
		}
		if got == nil || got.String() != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}
