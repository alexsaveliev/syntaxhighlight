package syntaxhighlight

var JavaLexer = &RegexpLexer{Rules: map[string][]RegexpRule {
	`root`: []RegexpRule{
			NewRule(`[^\S\n]+`, Text),
			NewRule(`//.*?\n`, Comment_Single),
			NewRule(`/\*.*?\*/`, Comment_Multiline),
			// keywords: go before method names to avoid lexing "throw new XYZ" as a method signature
			NewRule(`(assert|break|case|catch|continue|default|do|else|finally|for|if|goto|instanceof|new|return|switch|this|throw|try|while)\b`, Keyword),
			// method names
			NewActionRule(`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`,ByGroups(UsingThis(), Name_Function, Text, Operator)),
			NewRule(`@[^\W\d][\w.]*`, Name_Decorator),
			NewRule(`(abstract|const|enum|extends|final|implements|native|private|protected|public|static|strictfp|super|synchronized|throws|transient|volatile)\b`, Keyword_Declaration),
			NewRule(`(boolean|byte|char|double|float|int|long|short|void)\b`, Keyword_Type),
			NewActionRule(`(package)(\s+)`, ByGroups(Keyword_Namespace, Text), `import`),
			NewRule(`(true|false|null)\b`, Keyword_Constant),
			NewActionRule(`(class|interface)(\s+)`, ByGroups(Keyword_Declaration, Text), `class`),
			NewActionRule(`(import)(\s+)`, ByGroups(Keyword_Namespace, Text), `import`),
			NewRule(`"(\\\\|\\"|[^"])*"`, String),
			NewRule(`'\.'|'[^\]'|'\\u[0-9a-fA-F]{4}'`, String_Char),
			NewActionRule(`(\.)((?:[^\W\d]|\$)[\w$]*)`, ByGroups(Operator, Name_Attribute)),
			NewRule(`^\s*([^\W\d]|\$)[\w$]*:`, Name_Label),
			NewRule(`([^\W\d]|\$)[\w$]*`, Name),
			NewRule(`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator),
			NewRule(`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, Number_Float),
			NewRule(`0x[0-9a-fA-F]+`, Number_Hex),
			NewRule(`[0-9]+(_+[0-9]+)*L?`, Number_Integer),
			NewRule(`\n`, Text),
	},
	`class`: []RegexpRule{
		NewRule(`([^\W\d]|\$)[\w$]*`, Name_Class, `#pop`),
	},
	`import`: []RegexpRule{
		NewRule(`[\w.]+\*?`, Name_Namespace, `#pop`),
	},
}}