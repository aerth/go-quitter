package main

import (
	"fmt"
	"os"
)

func runtests() {
	testCutline()
	os.Exit(0)
}
func testCutline() {
	lines := cutline(10, "Operands denote the elementary values in an expression. An operand may be a literal, a (possibly qualified) non-blank identifier denoting a constant, variable, or function, a method expression yielding a function, or a parenthesized expression. ")
	for _, line := range lines {
		fmt.Println(line)
	}
}

var vtests2 = []struct {
	got  string
	gold string
}{
	{"1234567890223456789032 kjh kjh kjh kjh 34567890abcabcabcabcakdfjhskjh kjh flkjhsdflkjsdhflksjdhfslkjdfhslkjfhlkjsdhflksjdhf", "3234567890"},
}

var vtests = []struct {
	got  string
	gold string
}{
	{"lolzds kjs", "lolzds kjs"},
	{"lolzds kjshlkjhkljhlkjh", "lolzds kjs"},
	{"1234567890123456789", "1234567890"},
}
