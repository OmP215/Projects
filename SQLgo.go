package gosql

import (
	"fmt" //formatted io (ie printf, scanf)
	//"strings"//functions to manipulate strings
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

func lexKeyword(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing keywords
	return nil, cur, false
}

func lexSymbol(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing symbols
	return nil, cur, false
}

func lexString(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing strings
	return nil, cur, false
}

func lexNumeric(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing numeric values
	return nil, cur, false
}

func lexIdentifier(source string, cur cursor) (*token, cursor, bool) {
	// Implementation for lexing identifiers
	return nil, cur, false
}
