package syntaxhighlight

import (
	"reflect"
	"regexp"
)

// Matcher is a function that returns array of slices if source matches some conditions
// (an example is regexp.Regexp.FindSubmatch)
// source input byte slice
// returns array of slices where first element is full match and each next element is a sub-match or nil if source does not match conditions
type Matcher func(source []byte) [][]byte

// Produces lexer rules by compiling regular expressions using specific flags (such as DOTALL, MULTILINE and so on)
type FlagsRuleMaker struct {
	// flags to be used, for example `ims`
	flags string
}

// Produces new lexer rule that returns single token if source matches given RE at the current position
// - pattern - RE pattern to use (will be combined with flags and will start with ^)
// - ttype - token type to produce
// - state - optional state(s) to apply to lexer's stack
func (self FlagsRuleMaker) Token(pattern string, ttype *TokenType, state ...string) RegexpRule {
	return RegexpRule{matcher: self.makeRegexp(pattern), ttype: ttype, states: state}
}

// Produces new lexer rule that returns token(s) if source matches given RE at the current position. 
// Tokens are produced by a given function
// - pattern - RE pattern to use (will be combined with flags and will start with ^)
// - action - function that will return tokens based on current match
// - state - optional state(s) to apply to lexer's stack
func (self FlagsRuleMaker) Action(pattern string, action RuleAction, state ...string) RegexpRule {
	return RegexpRule{matcher: self.makeRegexp(pattern), action: action, states: state}
}

// Produces new lexer rule that returns single token if source matches given matcher object at the current position
// - matcher - matcher object
// - ttype - token type to produce
// - state - optional state(s) to apply to lexer's stack
func (self FlagsRuleMaker) MatcherToken(matcher Matcher, ttype *TokenType, state ...string) RegexpRule {
	return RegexpRule{matcher: matcher, ttype: ttype, states: state}
}

// Produces new lexer rule that returns token(s) if source matches given matcher object at the current position. 
// Tokens are produced by a given function
// - matcher - matcher object
// - action - function that will return tokens based on current match
// - state - optional state(s) to apply to lexer's stack
func (self FlagsRuleMaker) MatcherAction(matcher Matcher, action RuleAction, state ...string) RegexpRule {
	return RegexpRule{matcher: matcher, action: action, states: state}
}


// Adjust RE to include flags and force anchor mode
func (self FlagsRuleMaker) makeRegexp(pattern string) Matcher {
	p := pattern
	if self.flags != `` {
		p = `(?` + self.flags + `)` + p
	}
	return regexp.MustCompile(`^` + p).FindSubmatch
}

var (
	// Predefined rule maker (multiline, case-insensitive, dotall)
	MSI = FlagsRuleMaker{`msi`}
	// Predefined rule maker (multiline, dotall)
	MS = FlagsRuleMaker{`ms`}
)

// Function that may produce zero or more tokens based on a given matcher
// - lexer current lexer object
// - matches array of slices where first item is a full match and each next item is matcher's group
type RuleAction func(lexer Lexer, matches [][]byte) []Token

// Defines lexer rule
type RegexpRule struct {
	// matcher object that will identify if rule matches source at the current position
	matcher Matcher
	// if rule produces fixed token, holds token type to produce
	ttype *TokenType
	// optional list of states to be applied to lexer's stack 
	states []string
	// function to produce tokens instead of producing single one
	action RuleAction
}

// Either produces single token per each submatch (if argument is a TokenType) or calls custom functions to produce tokens from each submatch
// (if argument is a function)
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

// Produces tokens using current lexer from a substring of source code
func UsingThis() RuleAction {
	return func(lexer Lexer, matches [][]byte) []Token {
		return lexer.GetTokens(matches[0])
	}
}