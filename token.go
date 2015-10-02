package syntaxhighlight

import "fmt"

type Token struct {
	Text []byte
	Type *TokenType
}

func (self Token) String() string {
	return fmt.Sprintf("`%s` [%s]", string(self.Text), self.Type.Name())
}