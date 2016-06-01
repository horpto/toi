package main

import (
	"errors"

	//patched for AddHeaders
	"github.com/horpto/termtables"
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
	s := &Scheme{In: in, Out: out, Memory: memory}
	return s, nil
}

func (s *Scheme) String() (buffer string) {
	buffer += "input: " + s.In + "\n"
	buffer += "output: " + s.Out + "\n"

	memories := ""
	vars := ""
	for v, node := range s.Memory {
		memories += v + ","
		vars += v + ": " + node.String() + "\n"
	}
	if memories != "" {
		buffer += "memory: " + memories + "\n"
		buffer += vars
	}
	return buffer
}

func (s Scheme) calculate(namespace boolParser.Namespace) (boolParser.Namespace, error) {
	if _, ok := namespace[s.Out]; ok {
		return nil, errors.New("namespace contains output var: " + s.Out)
	}
	node, ok := s.Memory[s.Out]
	if !ok {
		return nil, errors.New("Had no formula for output var: " + s.Out)
	}
	out, err := node.Calculate(namespace)
	if err != nil {
		return nil, err
	}
	namespace[s.Out] = out

	memory := make(boolParser.Namespace, len(s.Memory))
	memory[s.Out] = out
	for x, n := range s.Memory {
		r, err := n.Calculate(namespace)
		if err != nil {
			delete(namespace, s.Out)
			return memory, err
		}
		memory[x] = r
	}
	delete(namespace, s.Out)
	return memory, nil
}

func (s Scheme) calculateOutputWord(signals []bool) ([]bool, error) {
	out := make([]bool, len(signals))
	namespace := make(map[string]bool, len(s.Memory))
	for k, _ := range s.Memory {
		namespace[k] = false
	}
	for i, signal := range signals {
		namespace[s.In] = signal
		delete(namespace, s.Out)

		res, err := s.calculate(namespace)
		if err != nil {
			return nil, err
		}
		out[i] = res[s.Out]
		namespace = res
	}
	return out, nil
}

func (s Scheme) createTruthTable() (*TruthTable, error) {
	// Фиксируем порядок имен переменных,
	// чтобы при итерации назначать новые значения
	memory := []string{}
	for k, _ := range s.Memory {
		if k == s.In || k == s.Out {
			continue
		}
		memory = append(memory, k)
	}
	inputs := append([]string{}, s.In)
	inputs = append(inputs, memory...)
	outputs := append([]string{}, s.Out)
	outputs = append(outputs, memory...)

	tableLength := 1 << uint32(len(s.Memory))
	vars := make([][]bool, tableLength)
	values := make([][]bool, tableLength)

	for i := 0; i < tableLength; i++ {
		namespace := make(map[string]bool, len(inputs))
		varsRow := make([]bool, len(inputs))
		c := i
		for j := len(inputs) - 1; j >= 0; j-- {
			v := c&1 == 1
			c = c >> 1
			namespace[inputs[j]] = v
			varsRow[j] = v
		}

		res, err := s.calculate(namespace)
		if err != nil {
			return nil, err
		}

		valuesRow := make([]bool, len(outputs))
		for i, v := range outputs {
			valuesRow[i] = res[v]
		}

		vars[i] = varsRow
		values[i] = valuesRow
	}
	tt := &TruthTable{
		inputs:  inputs,
		outputs: outputs,
		vars:    vars,
		values:  values,
	}
	return tt, nil
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
	table := termtables.CreateTable()
	table.SetModeTerminal()

	for _, v := range tt.inputs {
		table.AddHeaders(v)
	}
	for _, v := range tt.outputs {
		table.AddHeaders(v)
	}

	// assume len(tt.inputs) == tt.vars[0] and len(tt.output) == len(tt.outputs)
	// and actually len(tt.inputs) == len(tt.output)
	for i, vars := range tt.vars {
		row := table.AddRow()
		for _, v := range vars {
			row.AddCell(boolToString(v))
		}
		for _, v := range tt.values[i] {
			row.AddCell(boolToString(v))
		}
	}

	return table.Render()
}
