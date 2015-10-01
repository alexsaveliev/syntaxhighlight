package syntaxhighlight

type Formatter interface {
	Format(tokens []Token) string
}