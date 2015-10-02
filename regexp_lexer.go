package syntaxhighlight

import (
	"strconv"
)

// Lexer that uses mostly RE to produce tokens
type RegexpLexer struct {
	// map of rules to produce tokens in form state => []rule
	Rules map[string][]RegexpRule
}

// Produces tokens using RE
// source is a slice to produce tokens from
// stack (optional) is initial lexer's state, if omited then 'root' is assumed
// returns array of tokens found
func (self RegexpLexer) GetTokens(source []byte, stack ...[]string) []Token {
	pos := 0
	var statestack []string

	if len(stack) > 0 {
		statestack = stack[0]
	} else {
		statestack = []string{`root`}
	}
	ret := make([]Token, 0, 20)
	l := len(source)

	rules := self.Rules[statestack[len(statestack) - 1]]
	for {
   		if pos == l {
   			break
   		}
		match := false
		slice := source[pos:]
		for _, rule := range rules {
			matcher := rule.matcher(slice)
			if matcher == nil {
				continue
			}
			if rule.ttype != nil {
				ret = append(ret, 
					Token{Text: matcher[0], 
						Type: rule.ttype})
			} else {
				tokens := rule.action(self, matcher)
				ret = append(ret, tokens...)
			}
			match = true
			pos += len(matcher[0])
			statestack = updateStack(statestack[0:], rule)
			break
		}
		if !match {
			if source[pos] == '\n' {
				statestack = []string{`root`}
				rules = self.Rules[`root`]
				ret = append(ret, 
					Token{Text: source[pos:pos+1], 
						Type: Text})
				pos++
				continue
			}
   			ret = append(ret, 
   				Token{Text: source[pos:pos+1], 
   					Type: Error})
			pos += 1
		}
	  	rules = self.Rules[statestack[len(statestack) - 1]]
	}
	return ret
}

// updates stack state using given rule
// rule may define zero or more transitions to be applied, such as
// - #pop - pops stack 
// - #push - pushed head stack item on the top of stack 
// - <N> (number) - leave only N items in stack
// - <state> - put state to stack
func updateStack(stack []string, rule RegexpRule) []string {
	for _, state := range rule.states {
    	if state == `#pop` {
    		stack = stack[:len(stack) - 1]
    	} else if state == `#push` {
    		stack = append(stack, stack[len(stack) - 1])
    	} else if i, err := strconv.Atoi(state); err == nil {
    		stack = stack[:i]
    	} else {
    		stack = append(stack, state)
    	}
	}
	return stack
}