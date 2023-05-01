package token

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	FilePath string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 123456
	STRING = "STRING"
	FLOAT  = "FLOAT" // 123.45

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	EQ = "=="
	NE = "!="
	LT = "<"
	GT = ">"

	PIPE      = "|"
	AMPERSAND = "&"
	MOD       = "%"
	CARET     = "^"
	TILDE     = "~"

	//boolean constants
	TRUE  = "true"
	FALSE = "false"
	NONE  = "NONE"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	WHILE    = "WHILE"
	FOR      = "FOR"
	AND      = "AND"
	AS       = "AS"
	ASSERT   = "ASSERT"
	BREAK    = "BREAK"
	CLASS    = "CLASS"
	CONTINUE = "CONTINUE"
	DEF      = "DEF"
	DEL      = "DEL"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	EXCEPT   = "EXCEPT"
	EXEC     = "EXEC"
	FINALLY  = "FINALLY"
	FROM     = "FROM"
	GLOBAL   = "GLOBAL"
	IMPORT   = "IMPORT"
	IN       = "IN"
	IS       = "IS"
	LAMBDA   = "LAMBDA"
	NONLOCAL = "NONLOCAL"
	NOT      = "NOT"
	OR       = "OR"
	PASS     = "PASS"
	PRINT    = "PRINT"
	RAISE    = "RAISE"
	RETURN   = "RETURN"
	TRY      = "TRY"
	WITH     = "WITH"
	YIELD    = "YIELD"
	DEDENT   = "DEDENT"
	COLON    = ":"
)

var Keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"if":       IF,
	"while":    WHILE,
	"for":      FOR,
	"else":     ELSE,
	"return":   RETURN,
	"and":      AND,
	"as":       AS,
	"assert":   ASSERT,
	"break":    BREAK,
	"class":    CLASS,
	"continue": CONTINUE,
	"def":      DEF,
	"del":      DEL,
	"elif":     ELIF,
	"except":   EXCEPT,
	"exec":     EXEC,
	"finally":  FINALLY,
	"from":     FROM,
	"global":   GLOBAL,
	"import":   IMPORT,
	"in":       IN,
	"is":       IS,
	"lambda":   LAMBDA,
	"nonlocal": NONLOCAL,
	"not":      NOT,
	"or":       OR,
	"pass":     PASS,
	"print":    PRINT,
	"raise":    RAISE,
	"try":      TRY,
	"with":     WITH,
	"yield":    YIELD,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
