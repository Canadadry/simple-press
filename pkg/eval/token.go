package eval

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

const (
	ILLEGAL      TokenType = "ILLEGAL"
	EOF          TokenType = "EOF"
	NUMBER       TokenType = "NUMBER"
	LEFT_PAR     TokenType = "LEFT_PARENTHESIS"
	RIGHT_PAR    TokenType = "RIGHT_PARENTHESIS"
	PLUS         TokenType = "PLUS"
	MINUS        TokenType = "MINUS"
	STAR         TokenType = "STAR"
	SLASH        TokenType = "SLASH"
	DOUBLE_STAR  TokenType = "DOUBLE_STAR"
	DOUBLE_SLASH TokenType = "DOUBLE_SLASH"
	PERCENT      TokenType = "PERCENT"
)
