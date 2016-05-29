# toi


Грамматика:

main := statement
statement := intersection { [+ \ - ] statement}
intersection := expression { '\*' intersection }


//a + b * c == a + (b * c) // OK
//a + b * c == (a + b) * c // fail
b * c + a == (b * c) + a


expression := '!' expression | const | ident | (statement) | [statement]

const := binDig

ident := alph {alph | digit}*
binDig := '0' | '1'

digit := binDig | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
