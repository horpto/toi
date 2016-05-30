package main

import (
	"errors"

	"github.com/horpto/toi/lib"
)

type Node interface {
	boolParser.Node
}

type Scheme struct {
	In     string          // имя входного аргумента
	Out    string          // имя выходного аргумента
	Memory map[string]Node // словарь: имя переменной - как высчитывается.
}

func newScheme(in string, out string, mems map[string]string) (*Scheme, error) {
	if _, ok := mems[in]; ok {
		return nil, errors.New("Input variable '" + in + "' has formula")
	}
	if _, ok := mems[out]; !ok {
		return nil, errors.New("Output variable '" + out + "' has no formula")
	}

	memory := make(map[string]Node, len(mems))
	for k, v := range mems {
		mnode, err := boolParser.ParseString(v)
		if err != nil {
			return nil, err
		}
		memory[k] = mnode
	}
	s := Scheme{In: in, Out: out, Memory: memory}
	return &s, nil
}

func (s Scheme) calculate(namespace boolParser.Namespace) (bool, boolParser.Namespace, error) {
	if _, ok := namespace[s.Out]; ok {
		return false, nil, errors.New("namespace contains output var: " + s.Out)
	}
	node, ok := s.Memory[s.Out]
	if !ok {
		return false, nil, errors.New("Had no formula for output var: " + s.Out)
	}
	out, err := node.Calculate(namespace)
	if err != nil {
		return false, nil, err
	}
	namespace[s.Out] = out

	memory := make(boolParser.Namespace, len(s.Memory))
	memory[s.Out] = out
	for x, n := range s.Memory {
		r, err := n.Calculate(namespace)
		if err != nil {
			delete(namespace, s.Out)
			return out, memory, err
		}
		memory[x] = r
	}
	delete(namespace, s.Out)
	return out, memory, nil
}
