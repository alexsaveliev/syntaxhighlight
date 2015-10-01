package syntaxhighlight

type Lexer interface {
	GetTokens(source []byte, stack ...[]string) []Token
}