package syntaxhighlight

import "fmt"

type Token struct {
	Text []byte
	Offset int
	Type *TokenType
}

func (self Token) String() string {
	return fmt.Sprintf("`%s` at %d [%s]", string(self.Text), self.Offset, self.Type.Name())
}