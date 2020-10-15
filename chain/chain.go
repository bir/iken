package chain

import (
	"github.com/valyala/fasthttp"
)

// Constructor is the function contract used to make RequestHandlers.
type Constructor func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Chain enables easy definition of custom handler chains.
type Chain []Constructor

// New creates a new chain with the supplied Constructors.
func New(cc ...Constructor) Chain {
	return cc
}

// Handler compiles the chain to a single RequestHandler.
func (c Chain) Handler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	l := len(c) - 1
	// Execute in reverse order
	for i := range c {
		h = c[l-i](h)
	}

	return h
}

// Append adds the Constructor to the end of the execution chain.
func (c Chain) Append(cc ...Constructor) Chain {
	newChain := make(Chain, 0, len(c)+len(cc))
	newChain = append(newChain, c...)
	newChain = append(newChain, cc...)

	return newChain
}

// Prepend adds the Constructor to the start of the execution chain.
func (c Chain) Prepend(cc ...Constructor) Chain {
	newChain := make(Chain, 0, len(c)+len(cc))
	newChain = append(newChain, cc...)
	newChain = append(newChain, c...)

	return newChain
}
