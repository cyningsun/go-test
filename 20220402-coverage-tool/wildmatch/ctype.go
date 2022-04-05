package wildmatch

import "unicode"

const (
	GIT_SPACE          = 0x01
	GIT_DIGIT          = 0x02
	GIT_ALPHA          = 0x04
	GIT_GLOB_SPECIAL   = 0x08
	GIT_REGEX_SPECIAL  = 0x10
	GIT_PATHSPEC_MAGIC = 0x20
	GIT_CNTRL          = 0x40
	GIT_PUNCT          = 0x80
)

const (
	S = GIT_SPACE
	A = GIT_ALPHA
	D = GIT_DIGIT
	G = GIT_GLOB_SPECIAL   /* *, ?, [, \\ */
	R = GIT_REGEX_SPECIAL  /* $, (, ), +, ., ^, {, | */
	P = GIT_PATHSPEC_MAGIC /* other non-alnum, except for ] and } */
	X = GIT_CNTRL
	U = GIT_PUNCT
	Z = GIT_CNTRL | GIT_SPACE
)

var (
	saneCtype = [256]rune{
		X, X, X, X, X, X, X, X, X, Z, Z, X, X, Z, X, X, /*   0.. 15 */
		X, X, X, X, X, X, X, X, X, X, X, X, X, X, X, X, /*  16.. 31 */
		S, P, P, P, R, P, P, P, R, R, G, R, P, P, R, P, /*  32.. 47 */
		D, D, D, D, D, D, D, D, D, D, P, P, P, P, P, G, /*  48.. 63 */
		P, A, A, A, A, A, A, A, A, A, A, A, A, A, A, A, /*  64.. 79 */
		A, A, A, A, A, A, A, A, A, A, A, G, G, U, R, P, /*  80.. 95 */
		P, A, A, A, A, A, A, A, A, A, A, A, A, A, A, A, /*  96..111 */
		A, A, A, A, A, A, A, A, A, A, A, R, R, U, P, X, /* 112..127 */
		/* Nothing in the 128.. range */
	}
	hexvalTable = [256]rune{
		-1, -1, -1, -1, -1, -1, -1, -1, /* 00-07 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 08-0f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 10-17 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 18-1f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 20-27 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 28-2f */
		0, 1, 2, 3, 4, 5, 6, 7, /* 30-37 */
		8, 9, -1, -1, -1, -1, -1, -1, /* 38-3f */
		-1, 10, 11, 12, 13, 14, 15, -1, /* 40-47 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 48-4f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 50-57 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 58-5f */
		-1, 10, 11, 12, 13, 14, 15, -1, /* 60-67 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 68-67 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 70-77 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 78-7f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 80-87 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 88-8f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 90-97 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* 98-9f */
		-1, -1, -1, -1, -1, -1, -1, -1, /* a0-a7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* a8-af */
		-1, -1, -1, -1, -1, -1, -1, -1, /* b0-b7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* b8-bf */
		-1, -1, -1, -1, -1, -1, -1, -1, /* c0-c7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* c8-cf */
		-1, -1, -1, -1, -1, -1, -1, -1, /* d0-d7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* d8-df */
		-1, -1, -1, -1, -1, -1, -1, -1, /* e0-e7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* e8-ef */
		-1, -1, -1, -1, -1, -1, -1, -1, /* f0-f7 */
		-1, -1, -1, -1, -1, -1, -1, -1, /* f8-ff */
	}
)

func saneIstest(x rune, mask rune) bool {
	return (saneCtype[x] & mask) != 0
}

func isGlobSpecial(x rune) bool {
	return saneIstest(x, GIT_GLOB_SPECIAL)
}

func isascii(x rune) bool {
	return unicode.Is(unicode.ASCII_Hex_Digit, x)
}

func isspace(x rune) bool {
	return saneIstest(x, GIT_SPACE)
}
func isdigit(x rune) bool {
	return saneIstest(x, GIT_DIGIT)
}
func isalpha(x rune) bool {
	return saneIstest(x, GIT_ALPHA)

}
func isalnum(x rune) bool {
	return saneIstest(x, GIT_ALPHA|GIT_DIGIT)
}
func isprint(x rune) bool {
	return ((x) >= 0x20 && (x) <= 0x7e)
}
func islower(x rune) bool {
	return saneIscase(x, true)
}

func isupper(x rune) bool {
	return saneIscase(x, false)
}

func iscntrl(x rune) bool {
	return saneIstest(x, GIT_CNTRL)
}
func ispunct(x rune) bool {
	return unicode.IsPunct(x)
}

func isxdigit(x rune) bool {
	return hexvalTable[x] != -1
}
func tolower(x rune) rune {
	return saneCase(x, 0x20)
}
func toupper(x rune) rune {
	return saneCase(x, 0)
}
func isPathSpecMagic(x rune) bool {
	return saneIstest(x, GIT_PATHSPEC_MAGIC)
}

func saneCase(x, high rune) rune {
	if saneIstest(x, GIT_ALPHA) {
		x = (x & ^0x20) | high
	}

	return x
}

func saneIscase(x rune, isLower bool) bool {
	if !saneIstest(x, GIT_ALPHA) {
		return false
	}

	if isLower {
		return (x & 0x20) != 0
	}
	return (x & 0x20) == 0
}
