package syntaxhighlight

var JavaLexer = &RegexpLexer{Rules: map[string][]RegexpRule {
	`root`: []RegexpRule{
			MSI.Token(`[^\S\n]+`, Text),
			MSI.Token(`//.*?\n`, Comment_Single),
			MSI.Token(`/\*.*?\*/`, Comment_Multiline),
			// keywords: go before method names to avoid lexing "throw new XYZ" as a method signature
			MSI.Token(`(assert|break|case|catch|continue|default|do|else|finally|for|if|goto|instanceof|new|return|switch|this|throw|try|while)\b`, Keyword),
			// method names
			MSI.Action(`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`,ByGroups(UsingThis(), Name_Function, Text, Operator)),
			MSI.Token(`@[^\W\d][\w.]*`, Name_Decorator),
			MSI.Token(`(abstract|const|enum|extends|final|implements|native|private|protected|public|static|strictfp|super|synchronized|throws|transient|volatile)\b`, Keyword_Declaration),
			MSI.Token(`(boolean|byte|char|double|float|int|long|short|void)\b`, Keyword_Type),
			MSI.Action(`(package)(\s+)`, ByGroups(Keyword_Namespace, Text), `import`),
			MSI.Token(`(true|false|null)\b`, Keyword_Constant),
			MSI.Action(`(class|interface)(\s+)`, ByGroups(Keyword_Declaration, Text), `class`),
			MSI.Action(`(import)(\s+)`, ByGroups(Keyword_Namespace, Text), `import`),
			MSI.Token(`"(\\\\|\\"|[^"])*"`, String),
			MSI.Token(`'\.'|'[^\]'|'\\u[0-9a-fA-F]{4}'`, String_Char),
			MSI.Action(`(\.)((?:[^\W\d]|\$)[\w$]*)`, ByGroups(Operator, Name_Attribute)),
			MSI.Token(`^\s*([^\W\d]|\$)[\w$]*:`, Name_Label),
			MSI.Token(`([^\W\d]|\$)[\w$]*`, Name),
			MSI.Token(`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator),
			MSI.Token(`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, Number_Float),
			MSI.Token(`0x[0-9a-fA-F]+`, Number_Hex),
			MSI.Token(`[0-9]+(_+[0-9]+)*L?`, Number_Integer),
			MSI.Token(`\n`, Text),
	},
	`class`: []RegexpRule{
		MSI.Token(`([^\W\d]|\$)[\w$]*`, Name_Class, `#pop`),
	},
	`import`: []RegexpRule{
		MSI.Token(`[\w.]+\*?`, Name_Namespace, `#pop`),
	},
}}