package gosql

import (
	"fmt" //formatted io (ie printf, scanf)
	"strings"//functions to manipulate strings
)

type location struct {
	line uint
	col  uint
}

type keyword string

// Defines keywords for query
const (
	selectKeyword keyword = "select"
	fromKeyword   keyword = "from"
	asKeyword     keyword = "as"
	tableKeyword  keyword = "table"
	createKeyword keyword = "create"
	insertKeyword keyword = "insert"
	intoKeyword   keyword = "into"
	valuesKeyword keyword = "values"
	intKeyword    keyword = "int"
	textKeyword   keyword = "text"
)

type symbol string

// Defines symbols
const (
	semicolonSymbol        symbol = ";"
	asteriskSymbol         symbol = "*"
	commaSymbol            symbol = ","
	leftparenthesisSymbol  symbol = "("
	rightparenthesisSymbol symbol = ")"
)

type tokenKind uint

// tokens for parser
const (
	keywordKind tokenKind = iota
	symbolKind
	identifierKind
	stringKind
	numericKind
)

type token struct {
	value string
	kind  tokenKind
	loc   location
}

type cursor struct {
	pointer uint
	loc     location
}

func (t *token) equals(other *token) bool {
	return t.value == other.value && t.kind == other.kind
}

type lexer func(string, cursor) (*token, cursor, bool)

// main loop
func lex(source string) ([]*token, error) {
	tokens := []*token{}
	cur := cursor{}

lex:
	for cur.pointer < uint(len(source)) {
		lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
		for _, l := range lexers {
			if token, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				//omit null tokens for valid but empty syntax
				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.loc.line, cur.loc.col)
	}

	return tokens, nil
}
/* refer to, https://www.postgresql.org/docs/current/sql-syntax-lexical.html
*   for what constitutes a valid number
*/
func lexNumeric(source string, ic cursor) (*token, cursor, bool) {
	// Implementation for lexing numeric values
	cur := ic

	periodFound := false
	expMarkerFound := false

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		cur.loc.col++

		isDigit := c >= '0' c <= '9'
		isPeriod := c == '.'
		isExpMarker := c == 'e'

		//must start w/ a digit or period
		if cur.pointer == ic.pointer {
			if !isdigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}
		if isPeriod {
			if periodFound {
				return nul, ic, false
			}

			periodFound = true
			continue
		}

		if isExpMarker {
			if expMarkerFound {
				return nil, ic, false
			}

			//NO periods allowed after expMarker

			periodFound - true
			expMarkerFound = true

			// expMarker must be followed by digits
			if cur.pointer == uint(len(source)- 1) {
				return nil, ic, false
			}

			cNext := source[cur.pointer+1]
			if cNext == '-' || cNext == '+' {
				cur.pointer++
				cur.loc.col++
			}
			continue
		}
		if !isDigit {
			break
		}
	}
	//no chars accumulated
	if cur.pointer == ic.pointer {
		return nil, ic, false
	}
	return &token {
		value: source[ic.pointer:cur.pointer],
		loc: ic.loc,
		kind: numericKind
	}, cur, true
}

/*
Strings must start and end with a single apostophe. Can contain 
a single apostophe if its followed by another single postrophe.
Helper function for char delimited lexing logic to analyze identifiers
*/
func lexCharacterDelimited(source string, ic cursor, delimiter byte) (*token, cursor, bool) {
	// Implementation for lexing keywords
	cue := ic

	if len(source[cur.pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.pointer] != demiliter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.pointer++

	var value []byte
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter {
			// SQL escapes are w/ double chars, not backslash
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				return &token{
					value: string(value)
					loc: ic.loc
					kind: stringKind,
				}, cur, true
			} else {
				value = append(value, delimiter)
				cur.pointer++
				cur.loc.col++
			}
		}

		value = append(value, c)
		cur.loc.col++
	}
	return nil, ic, false
}

func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	// Implementation for lexing symbols
	c := source[ic.pointer]
	
	cur.pointer++
	cur.loc.col++
	
	switch c {
		//syntax that should be thrown away
	case '\n':
		cur.col.line++
		cur.col.loc = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	//syntax that should be kept
	symbols := []symbol{
		commaSymbol,
		leftparenthesisSymbol,
		rightparenthesisSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options,string(s))
	}

	match := longestMatch(source, ic, options)
	//unknown char
	if match == '' {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		value: match,
		loc: 	ic.loc,
		kind: symbolKind,
	}, cur, true
}

func lexString(source string, ic cursor) (*token, cursor, bool) {
	// Implementation for lexing strings
	return lexCharacterDelimited(source, ic '\'')
}

func lexKeyword(source string, ic cursor) (*token, cursor, bool) {
	cur := ic
	keywords := []keyword{
		selectKeyword,
		insertKeyword,
		valuesKeyword,
		tableKeyword,
		createKeyword,
		whereKeyword,
		fromKeyword,
		intoKeyword,
		textKeyword,
	}
	var options []string
	for _, k := range keywords{
		options = append(options, string(k))
	}
	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		value: match,
		kind: kind,
		loc: ic.loc,
	}, cur, true
}


//longestMatch iterates through a source string starting at the given cursor
//to find the longest matching among the provided options
func longestMatch(source string, ic cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := cur

	for cur.pointer < uint(len(source)) {
		value = append(value, strings.ToLower(string(source[cur.pointer]))...)
		cur.pointer++

		match: 
			for i, option := range options {
				for _, skip := range skipList {
					if i == skip {
					continue match
					}
				}
				//Deals with int vs into
				if option == string(value) {
					skipList == append(skipList, i)
					if len(options) > len(match) {
						match = option
					}
					continue
				}
				sharesPrefix := string(value) == options[:cur.pointer-ic.pointer]
				tooLong := len(value) > len(option)
				if tooLong || !sharesPrefix {
					skipList =append(skipList, i)
				}
			}
			if len(skipList) == len(options) {
				break
			}
	}
	return match
}


func lexIdentifier(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing identifiers
	return nil, cur, false
}
