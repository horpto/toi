package boolParser

type TokenStream interface {
	topToken() Token
	popToken()
	putToken(Token)
}

type ArrayTokenStream struct {
	eof   bool
	top   Token // cache
	stack []*Token
	ch    <-chan Token
}

func NewArrayTokenStream(ch <-chan Token) ArrayTokenStream {
	top, ok := <-ch
	if !ok {
		top = Token{_type: tokEOF, offset: 0}
	}
	return ArrayTokenStream{
		eof:   !ok,
		top:   top,
		stack: nil, // maybe fail with len()
		ch:    ch,
	}
}

func (ts ArrayTokenStream) topToken() Token {
	return ts.top
}

func (ts *ArrayTokenStream) popToken() {
	if ts.eof {
		return
	}
	if len(ts.stack) > 0 {
		lastItemIndex := len(ts.stack) - 1
		ts.top = *ts.stack[lastItemIndex]
		ts.eof = ts.top._type != tokEOF
		ts.stack = ts.stack[:lastItemIndex]
		return
	}

	top, ok := <-ts.ch
	if !ok {
		top = Token{_type: tokEOF, offset: 0}
	}
	ts.top = top
	ts.eof = ts.top._type == tokEOF
}

func (ts *ArrayTokenStream) putToken(token Token) {
	ts.stack = append(ts.stack, &ts.top)
	ts.top = token
}
