package conf

type Token int

const (
	invalid Token = iota

	literal_begin
	IDENT
	INT
	FLOAT
	IMAG
	STRING
	BOOL
	literal_end

	LBRACK // [
	LBRACE // {
	RBRACK // ]
	RBRACE // }

	keyword_begin
	INCLUDE // include directive
	keyword_end
)

var tokens = [...]string{
	IDENT:   "IDENT",
	INT:     "INT",
	FLOAT:   "FLOAT",
	IMAG:    "IMAG",
	STRING:  "STRING",
	BOOL:    "BOOL",
	LBRACK:  "[",
	LBRACE:  "{",
	RBRACK:  "]",
	RBRACE:  "}",
	INCLUDE: "include",
}

var keywords = map[string]Token{}

func init() {
	for i := keyword_begin + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func lookupToken(lit string) Token {
	if tok, ok := keywords[lit]; ok {
		return tok
	}
	return IDENT
}
