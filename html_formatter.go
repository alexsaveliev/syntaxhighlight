package syntaxhighlight

import (
	"fmt"
	"html"
	"strings"
)

type HtmlFormatter struct {
}

func (self HtmlFormatter) Format(tokensource []Token) []string {
	source := self.formatLines(tokensource)
		// if self.hl_lines {
		// 	source = self._highlight_lines(source)
		// } 
		//if not self.nowrap {
			// if self.linenos == 2 {
			// 	source = self._wrap_inlinelinenos(source)
			// }
			// if self.lineanchors {
			// 	source = self._wrap_lineanchors(source)
			// }
			//if self.linespans {
			// 	source = self._wrap_linespans(source)
			// }
			//source = self.wrap(source, outfile)
			// if self.linenos == 1 {
			// 	source = self._wrap_tablelinenos(source)
			// }
			//if self.full {
			//	source = self._wrap_full(source, outfile)
			//}
		//}
	return source
}

func (self HtmlFormatter) formatLines(tokensource []Token) []string {

	// TODO
	lsep := "\n"

	lspan := ``
	line := ``
	ret := make([]string, 0)

	for _, tok := range tokensource {
		ttype := tok.Type
		value := tok.Text

		cls := self.getCssClass(ttype)
		var cspan string
		if cls != `` {
			cspan = fmt.Sprintf(`<span class="%s">`, cls)
		} else {
			cspan = ``
		}

		parts := strings.Split(html.EscapeString(string(value)), "\n")
		len := len(parts)

		// for all but the last line
		for i := 0; i < len - 1; i++ {
			part := parts[i]
			if line != `` {
				if lspan != cspan {
					line += stringIf(lspan, `</span>`) + cspan + part + stringIf(cspan, `</span>`) + lsep
				} else { 
					// both are the same
					line += part + stringIf(lspan, `</span>`) + lsep
				}
				ret = append(ret, line)
				line = ``
			} else if part != `` {
				ret = append(ret, cspan + part + stringIf(cspan, `</span>`) + lsep)
			} else {
				ret = append(ret, lsep)
			}
		}
		// for the last line
		if line != `` && parts[len - 1] != `` {
			if lspan != cspan {
				line += stringIf(lspan, `</span>`) + cspan + parts[len - 1]
				lspan = cspan
			} else {
				line += parts[len - 1]
			}
		} else if parts[len - 1] != `` {
			line = cspan + parts[len - 1]
			lspan = cspan
		}
		// else we neither have to open a new span nor set lspan
	}

	if line != `` {
		ret = append(ret, line + stringIf(lspan, `</span>`) + lsep)
	}
	return ret
}

func (self HtmlFormatter) getCssClass(ttype *TokenType) string {
	ttypeclass := getTtypeClass(ttype)
	if ttypeclass != `` {
		// TODO
		//return self.classprefix + ttypeclass
		return ttypeclass
	}
	return ``
}

func getTtypeClass(ttype *TokenType) string {
	if STANDARD_TYPES[ttype.Name()] != nil {
		return ttype.Name()
	}
	tokens := make([]string, 0)
	for ttype != nil {
		tokens = append([]string{ttype.Name()}, tokens...)
		ttype = ttype.parent
	} 
	return strings.Join(tokens, `-`)
}

func stringIf(in string, out string) string {
	if in == `` {
		return ``
	} else {
		return out
	}
}