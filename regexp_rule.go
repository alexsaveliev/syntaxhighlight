package syntaxhighlight

import (
	"reflect"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Matcher interface {
	Group(index int) []byte
}

type pseudoMatcher struct {
	data []byte
}

func (self pseudoMatcher) Group(index int) []byte {
	return self.data
}

type RuleAction func(lexer Lexer, start int, matcher Matcher) []Token

type RegexpRule struct {
	pattern pcre.Regexp
	ptext string
	ttype *TokenType
	states []string
	action RuleAction
}

func NewRule(pattern string, ttype *TokenType, state ...string) RegexpRule {
	return RegexpRule{pattern: matcher(p(pattern)), ptext: p(pattern), ttype: ttype, states: state}
}

func NewActionRule(pattern string, action RuleAction, state ...string) RegexpRule {
	return RegexpRule{pattern: matcher(p(pattern)), ptext: p(pattern), action: action, states: state}
}

func matcher(pattern string) pcre.Regexp {
	return pcre.MustCompile(`^` + pattern, pcre.DOTALL | pcre.MULTILINE)
}

func p(pattern string) string {
	return `^` + pattern
}

func ByGroups(args ...interface{}) RuleAction {
	return func(lexer Lexer, start int, matcher Matcher) []Token {
		l := len(args)
		ret := make([]Token, 0, l)
		for i := 0; i < l; i++ {
			var t Token
			arg := args[i]
			vf := reflect.ValueOf(arg)
  			ftype := vf.Type()
  			if ftype.Kind() == reflect.Func {
  				tokens := arg.(RuleAction)(lexer, 0, pseudoMatcher{data: matcher.Group(i + 1)})
  				for _, tok := range tokens {
  					ret = append(ret, tok)
  				}
  			} else {
				t = Token{Text: matcher.Group(i + 1), Type: arg.(*TokenType)}		
				ret = append(ret, t)
  			}
		}
		return ret
	}
}

func UsingThis() RuleAction {
	return func(lexer Lexer, start int, matcher Matcher) []Token {
		return lexer.GetTokens(matcher.Group(0))
	}
}