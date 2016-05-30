package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func parseHeader(line, sep string) (string, error) {
	parts := strings.Split(line, sep)

	if len(parts) != 2 {
		return "", errors.New("Fail to parse line:" + line)
	}
	return strings.TrimSpace(parts[1]), nil
}

func createSchemeFromFile(fileName string) (*Scheme, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\r\n")
	exprs := make(map[string]string, len(lines))
	memory := []string{}

	var in, out string
	for i, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		switch {
		case strings.HasPrefix(line, "#"):
			continue
		case strings.HasPrefix(line, "input:") && in == "":
			in, err = parseHeader(line, ":")
			if err != nil {
				fmt.Errorf("miss line: %q", err.Error())
			}
		case strings.HasPrefix(line, "output:") && out == "":
			out, err = parseHeader(line, ":")
			if err != nil {
				fmt.Errorf("miss line: %q", err.Error())
				continue
			}
		case strings.HasPrefix(line, "memory:"):
			mem, err := parseHeader(line, ":")
			if err != nil {
				fmt.Errorf("miss line: %q", err.Error())
				continue
			}
			vars := strings.Split(mem, ",")
			for i := range vars {
				vars[i] = strings.TrimSpace(vars[i])
			}
			memory = append(memory, vars...)
		default:
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				fmt.Errorf("cannot parse line: %s", line)
				continue
			}
			exprs[parts[0]] = parts[1]
		}
		fmt.Printf("line %d: %s\r\n", i, line)
	}
	memory = append(memory, out)
	for _, mem := range memory {
		if _, ok := exprs[mem]; !ok {
			return nil, errors.New("Formula for '" + mem + "' not defined")
		}
	}
	return newScheme(in, out, exprs)
}

func main() {
	fileName := "./input.txt"

	s, err := createSchemeFromFile(fileName)
	fmt.Printf("%s", s, err)
}
