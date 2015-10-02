package syntaxhighlight

import (
	"reflect"
	"regexp"
)

type FlagsRuleMaker struct {
	flags string
}

func (self FlagsRuleMaker) Token(pattern string, ttype *TokenType, state ...string) RegexpRule {
	return RegexpRule{pattern: self.makeRegexp(pattern), ttype: ttype, states: state}
}

func (self FlagsRuleMaker) Action(pattern string, action RuleAction, state ...string) RegexpRule {
	return RegexpRule{pattern: self.makeRegexp(pattern), action: action, states: state}
}

func (self FlagsRuleMaker) makeRegexp(pattern string) *regexp.Regexp {
	p := pattern
	if self.flags != `` {
		p = `(?` + self.flags + `)` + p
	}
	return regexp.MustCompile(`^` + p)
}

var (
	MSI = FlagsRuleMaker{`msi`}
)


type RuleAction func(lexer Lexer, matches [][]byte) []Token

type RegexpRule struct {
	pattern *regexp.Regexp
	ttype *TokenType
	states []string
	action RuleAction
}

func ByGroups(args ...interface{}) RuleAction {
	return func(lexer Lexer, matches [][]byte) []Token {
		l := len(args)
		ret := make([]Token, 0, l)
		for i := 0; i < l; i++ {
			var t Token
			arg := args[i]
			vf := reflect.ValueOf(arg)
  			ftype := vf.Type()
  			if ftype.Kind() == reflect.Func {
  				ret = append(ret, arg.(RuleAction)(lexer, [][]byte {matches[i + 1]})...)
  			} else {
				t = Token{Text: matches[i + 1], Type: arg.(*TokenType)}		
				ret = append(ret, t)
  			}
		}
		return ret
	}
}

func UsingThis() RuleAction {
	return func(lexer Lexer, matches [][]byte) []Token {
		return lexer.GetTokens(matches[0])
	}
}