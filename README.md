# toi

## Задание:
разработать конструктор конечных автоматов, решающих  задачи анализа -
по составленной из заготовок логических элементов и задержек (не более 2-х) схеме с одним входом и одним выходом определяется таблица значений и вычисляется выходное слово по входному


### Грамматика:

```
main := statement

statement := intersection { [+ \ - ] statement}
intersection := expression { '\*' intersection }

expression := '!' expression | const | ident | (statement) | [statement]

const := binDig

ident := alph {alph | digit}*
binDig := '0' | '1'

digit := binDig | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
```

//a + b * c == a + (b * c) // OK
//a + b * c == (a + b) * c // fail
b * c + a == (b * c) + a

