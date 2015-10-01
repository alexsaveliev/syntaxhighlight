package syntaxhighlight

import (
	"reflect"
	"regexp"
)

type RuleAction func(lexer Lexer, source []byte, start int, matches []int) []Token

type RegexpRule struct {
	pattern *regexp.Regexp
	ttype *TokenType
	states []string
	action RuleAction
}

func NewRule(pattern string, ttype *TokenType, state ...string) RegexpRule {
	return RegexpRule{pattern: matcher(pattern), ttype: ttype, states: state}
}

func NewActionRule(pattern string, action RuleAction, state ...string) RegexpRule {
	return RegexpRule{pattern: matcher(pattern), action: action, states: state}
}

func matcher(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`^(?msi)` + pattern)
}

func ByGroups(args ...interface{}) RuleAction {
	return func(lexer Lexer, source []byte, start int, matches []int) []Token {
		l := len(args)
		ret := make([]Token, 0, l)
		for i := 0; i < l; i++ {
			var t Token
			index := (i + 1) * 2
			arg := args[i]
			vf := reflect.ValueOf(arg)
  			ftype := vf.Type()
  			if ftype.Kind() == reflect.Func {
  				tokens := arg.(RuleAction)(lexer, source[matches[index]:matches[index + 1]], 0, matches)
  				for _, tok := range tokens {
  					tok.Offset += start
  					ret = append(ret, tok)
  				}
  			} else {
				t = Token{Text: source[matches[index]:matches[index + 1]], Offset: start + matches[index], Type: arg.(*TokenType)}		
				ret = append(ret, t)
  			}
		}
		return ret
	}
}

func UsingThis() RuleAction {
	return func(lexer Lexer, source []byte, start int, matches []int) []Token {
		return lexer.GetTokens(source)
	}
}