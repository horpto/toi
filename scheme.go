package main

import (
	"errors"
	"strings"

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

type TruthTable struct {
	inputs  []string
	outputs []string
	vars    [][]bool
	values  [][]bool
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (tt *TruthTable) String() string {
	header := "|"
	paddings := make([]int, len(tt.inputs)+len(tt.outputs))
	l := 0
	for i, v := range tt.inputs {
		l += len(v)
		header += v + "|"
		paddings[i] = len(v) - 1
	}
	for i, v := range tt.outputs {
		l += len(v)
		header += v + "|"
		paddings[len(tt.inputs)+i] = len(v) - 1
	}
	header += "|\n"

	// body
	body := ""
	for i, vars := range tt.vars {
		line := "|"
		for j, v := range vars {
			line += strings.Repeat(" ", paddings[j]) + boolToString(v) + "|"
		}
		for j, v := range tt.values[i] {
			line += strings.Repeat(" ", paddings[len(vars)+j]) + boolToString(v) + "|"
		}
		line += "|\n"
		body += line
	}

	return header + body
}
