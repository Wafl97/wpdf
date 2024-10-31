package wpdf

func IsWhiteSpace(b byte) bool {
	return b == '\n' || b == '\r' || b == ' ' || b == '\t'
}

func IsDelimiter(b byte) bool {
	return b == '/' || b == '[' || b == ']' || b == '<' || b == '>' || b == '(' || b == ')'
}

func IsNumber(b byte) bool {
	return b == '0' ||
		b == '1' ||
		b == '2' ||
		b == '3' ||
		b == '4' ||
		b == '5' ||
		b == '6' ||
		b == '7' ||
		b == '8' ||
		b == '9'
}

func IsNumerical(b byte) bool {
	return b == '0' ||
		b == '1' ||
		b == '2' ||
		b == '3' ||
		b == '4' ||
		b == '5' ||
		b == '6' ||
		b == '7' ||
		b == '8' ||
		b == '9' ||
		b == '+' ||
		b == '-' ||
		b == '.'
}
