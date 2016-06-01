package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dixonwille/wmenu"
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
				if v := strings.TrimSpace(vars[i]); v != "" {
					vars[i] = v
				}
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

func ask(prompt string) string {
	answer := ""
	fmt.Print(prompt)
	fmt.Scanln(&answer)
	return strings.TrimSpace(answer)
}

func createSchemeFromStdin() (*Scheme, error) {
	in := ask("Введите имя входного параметра(x по умолчанию):")
	if in == "" {
		in = "x"
	}
	out := ask("Введите имя выходного параметра(y по умолчанию):")
	if out == "" {
		out = "y"
	}
	m := ask("Введите через запятую имена задержек:")
	vars := map[string]string{}
	for _, v := range strings.Split(m, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		vars[v] = ask("Введите лог.выражение для задержки '" + v + "'")
	}
	if _, ok := vars[out]; !ok {
		vars[out] = ask("Введите лог.выражение для выходного параметра")
	}
	return newScheme(in, out, vars)
}

func main() {
	var s *Scheme = nil

	exited := false
	for !exited {
		menu := wmenu.NewMenu("Выберите что-нибудь:")
		menu.Option("Выход", false, func() error {
			exited = true
			return nil
		})
		menu.Option("Ввести новую схему", false, func() error {
			s1, err := createSchemeFromStdin()
			if s1 != nil {
				s = s1
			}
			return err
		})
		menu.Option("Ввести новую схему из файла", false, func() error {
			fileName := ask("Введите путь до файла:")
			s1, err := createSchemeFromFile(fileName)
			if s1 != nil {
				s = s1
			}
			return err
		})
		menu.Option("Сохранить схему в файл", false, func() error {
			if s == nil {
				return errors.New("Введите сначала схему")
			}
			fileName := ask("Введите путь до файла:")
			var err error
			if fileName != "" {
				ioutil.WriteFile(fileName, []byte(s.String()), os.ModeType)
			}
			return err
		})
		menu.Option("Вывести таблицу истинности", false, func() error {
			if s == nil {
				return errors.New("Введите сначала схему")
			}
			tt, err := s.createTruthTable()
			if err != nil {
				return err
			}
			fmt.Print(tt.String())
			return nil
		})
		menu.Option("Найти выходное слово по входному", false, func() error {
			if s == nil {
				return errors.New("Введите сначала схему")
			}
			word := ask("Введите входное слово:")
			inputWord := []bool{}
			filteredWord := ""
			for _, w := range word {
				if w == '0' || w == '1' {
					inputWord = append(inputWord, w == '1')
					filteredWord += string(w)
				}
			}

			outputWord, err := s.calculateOutputWord(inputWord)
			if err != nil {
				return err
			}
			output := ""
			for _, w := range outputWord {
				if w {
					output += "1"
				} else {
					output += "0"
				}
			}
			fmt.Println(filteredWord)
			fmt.Println(output)
			return nil
		})
		menu.Run()
	}

}
