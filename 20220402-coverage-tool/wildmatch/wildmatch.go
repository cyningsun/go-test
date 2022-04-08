package wildmatch

import (
	"unicode"
)

const (
	NEGATE_CLASS  = '!'
	NEGATE_CLASS2 = '^'

	WM_CASEFOLD = 1
	WM_PATHNAME = 2

	WM_NOMATCH           = 1
	WM_MATCH             = 0
	WM_ABORT_ALL         = -1
	WM_ABORT_TO_STARSTAR = -2
)

/* *, ?, [, \\ */
func isGlobSpecial(x rune) bool {
	switch x {
	case '*', '?', '[', '\\':
		return true
	default:
		return false
	}
}

func equals(class, litmatch []rune) bool {
	if len(class) != len(litmatch) {
		return false
	}
	for i := 0; i < len(class); i++ {
		if class[i] != litmatch[i] {
			return false
		}
	}
	return true
}

// https://stackoverflow.com/questions/15767863/whats-the-difference-between-space-and-blank
// [:blank:]
//
//   Blank characters: space and tab.
//
// [:space:]
//
//     Space characters: in the 'C' locale, this is tab, newline,
//     vertical tab, form feed, carriage return, and space.
//
func IsBlank(c rune) bool {
	return c == ' ' || (c) == '\t'
}

func IsAlphaNum(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c)
}

func strchr(text []rune, c rune) int {
	for idx, each := range text {
		if each == c {
			return idx
		}
	}
	return len(text)
}

func doWild(pattern, text []rune, flags int) int {
	iText, iPattern := 0, 0
	for ; iPattern < len(pattern); iPattern, iText = iPattern+1, iText+1 {
		var (
			matched, matchSlash, negated int
			tCh, pCh, prevCh             rune
		)

		if iPattern < len(pattern) {
			pCh = pattern[iPattern]
		}

		if iText < len(text) {
			tCh = text[iText]
		}

		if iText == len(text) && pCh != rune('*') {
			return WM_ABORT_ALL
		}
		if (flags&WM_CASEFOLD) != 0 && unicode.IsUpper(tCh) {
			tCh = unicode.ToLower(tCh)
		}
		if (flags&WM_CASEFOLD) != 0 && unicode.IsUpper(pCh) {
			pCh = unicode.ToLower(pCh)
		}

		switch pCh {
		case '\\':
			/* Literal match with following character.  Note that the test
			 * in "default" handles the p[1] == '\0' failure case. */
			iPattern++
			pCh = pattern[iPattern]
			/* FALLTHROUGH */
		default:
			if tCh != pCh {
				return WM_NOMATCH
			}

			continue
		case '?':
			/* Match anything but '/'. */
			if (flags&WM_PATHNAME) != 0 && tCh == '/' {
				return WM_NOMATCH
			}

			continue
		case '*':
			iPattern++
			if iPattern < len(pattern) && pattern[iPattern] == '*' {
				prevP := rune(0)
				if iPattern-2 > 0 {
					prevP = pattern[iPattern-2]
				}

				for iPattern < len(pattern) && pattern[iPattern] == '*' {
					iPattern++
				}
				if (flags & WM_PATHNAME) == 0 {
					/* without WM_PATHNAME, '*' == '**' */
					matchSlash = 1
				} else if ((iPattern < len(pattern) && prevP < pattern[iPattern]) || prevP == '/') &&
					(iPattern == len(pattern) || pattern[iPattern] == '/' ||
						(pattern[iPattern] == '\\' && pattern[iPattern+1] == '/')) {
					/*
					 * Assuming we already match 'foo/' and are at
					 * <star star slash>, just assume it matches
					 * nothing and go ahead match the rest of the
					 * pattern with the remaining string. This
					 * helps make foo/<*><*>/bar (<> because
					 * otherwise it breaks C comment syntax) match
					 * both foo/bar and foo/a/bar.
					 */
					if iPattern < len(pattern) && pattern[iPattern] == '/' && doWild(pattern[iPattern+1:], text[iText:], flags) == WM_MATCH {
						return WM_MATCH
					}

					matchSlash = 1
				} else { /* WM_PATHNAME is set */
					matchSlash = 0
				}
			} else {
				/* without WM_PATHNAME, '*' == '**' */
				if (flags & WM_PATHNAME) != 0 {
					matchSlash = 0
				} else {
					matchSlash = 1
				}
			}

			if iPattern == len(pattern) {
				/* Trailing "**" matches everything.  Trailing "*" matches
				 * only if there are no more slash characters. */
				if matchSlash == 0 {
					if strchr(text[iText:], '/') != len(text[iText:]) {
						return WM_NOMATCH
					}
				}
				return WM_MATCH
			} else if matchSlash == 0 && pattern[iPattern] == '/' {
				/*
				 * _one_ asterisk followed by a slash
				 * with WM_PATHNAME matches the next
				 * directory
				 */
				idx := strchr(text[iText:], '/')
				if idx == len(text[iText:]) {
					return WM_NOMATCH
				}

				text = text[idx:]
				/* the slash is consumed by the top-level for loop */
				break
			}

			for {
				if iText == len(text) {
					break
				}
				/*
				 * Try to advance faster when an asterisk is
				 * followed by a literal. We know in this case
				 * that the string before the literal
				 * must belong to "*".
				 * If match_slash is false, do not look past
				 * the first slash as it cannot belong to '*'.
				 */
				if !isGlobSpecial(pattern[iPattern]) {
					pCh = pattern[iPattern]

					if (flags&WM_CASEFOLD) != 0 && unicode.IsUpper(pCh) {
						pCh = unicode.ToLower(pCh)
					}

					for iText < len(text) {
						tCh = text[iText]
						if matchSlash == 0 && tCh == '/' {
							break
						}
						if (flags&WM_CASEFOLD) != 0 && unicode.IsUpper(tCh) {
							tCh = unicode.ToLower(tCh)
						}
						if tCh == pCh {
							break
						}
						iText++
					}
					if tCh != pCh {
						return WM_NOMATCH
					}
				}
				if matched = doWild(pattern[iPattern:], text[iText:], flags); matched != WM_NOMATCH {
					if matchSlash != 0 || matched != WM_ABORT_TO_STARSTAR {
						return matched
					}
				} else if matchSlash == 0 && tCh == '/' {
					return WM_ABORT_TO_STARSTAR
				}
				iText++

				if iText < len(text) {
					tCh = text[iText]
				}
			}
			return WM_ABORT_ALL
		case '[':
			iPattern++
			if iPattern < len(pattern) {
				pCh = pattern[iPattern]
			}
			if pCh == NEGATE_CLASS2 {
				pCh = NEGATE_CLASS
			}

			/* Assign literal 1/0 because of "matched" comparison. */
			if pCh == NEGATE_CLASS {
				negated = 1
			} else {
				negated = 0
			}

			if negated != 0 { /* Inverted character class. */
				iPattern++
				if iPattern < len(pattern) {
					pCh = pattern[iPattern]
				}
			}

			prevCh = 0
			matched = 0
			for {
				if iPattern == len(pattern) {
					return WM_ABORT_ALL
				}
				if pCh == '\\' {
					iPattern++
					pCh = pattern[iPattern]
					if iPattern == len(pattern) {
						return WM_ABORT_ALL
					}
					if tCh == pCh {
						matched = 1
					}
				} else if pCh == '-' && prevCh != 0 && (iPattern+1) < len(pattern) && pattern[iPattern+1] != ']' {
					iPattern++
					pCh = pattern[iPattern]

					if pCh == '\\' {
						iPattern++
						pCh = pattern[iPattern]

						if iPattern == len(pattern) {
							return WM_ABORT_ALL
						}
					}
					if tCh <= pCh && tCh >= prevCh {
						matched = 1
					} else if (flags&WM_CASEFOLD) != 0 && unicode.IsLower(tCh) {
						tChUpper := unicode.ToUpper(tCh)
						if tChUpper <= pCh && tChUpper >= prevCh {
							matched = 1
						}
					}
					pCh = 0 /* This makes "prev_ch" get set to 0. */
				} else if pCh == '[' && pattern[iPattern+1] == ':' {
					var (
						s  []rune
						is int
						i  int
					)
					is = iPattern + 2
					s = pattern[is:]
					for {
						if iPattern == len(pattern) {
							break
						}

						pCh = pattern[iPattern]
						if pCh == ']' {
							break
						}
						iPattern++
					}
					if iPattern == len(pattern) {
						return WM_ABORT_ALL
					}

					i = iPattern - is - 1
					if i < 0 || pattern[iPattern-1] != ':' {
						pattern = pattern[is-2:]
						pCh = '['
						if tCh == pCh {
							matched = 1
						}
						continue
					}
					if equals(s[:i], []rune("alnum")) {
						if IsAlphaNum(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("alpha")) {
						if unicode.IsLetter(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("blank")) {
						if IsBlank(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("cntrl")) {
						if unicode.IsControl(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("digit")) {
						if unicode.IsDigit(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("graph")) {
						if unicode.IsGraphic(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("lower")) {
						if unicode.IsLower(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("print")) {
						if unicode.IsPrint(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("punct")) {
						if unicode.IsPunct(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("space")) {
						if unicode.IsSpace(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("upper")) {
						if unicode.IsUpper(tCh) {
							matched = 1
						} else if (flags&WM_CASEFOLD) != 0 && unicode.IsLower(tCh) {
							matched = 1
						}
					} else if equals(s[:i], []rune("xdigit")) {
						if unicode.Is(unicode.Hex_Digit, tCh) {
							matched = 1
						}
					} else { /* malformed [:class:] string */
						return WM_ABORT_ALL
					}
					pCh = 0 /* This makes "prev_ch" get set to 0. */
				} else if tCh == pCh {
					matched = 1
				}

				prevCh = pCh
				iPattern++
				if iPattern < len(pattern) {
					pCh = pattern[iPattern]
				}

				if pCh == ']' {
					break
				}
			}
			if matched == negated || (flags&WM_PATHNAME) != 0 && tCh == '/' {
				return WM_NOMATCH
			}

			continue
		}
	}

	if iText != len(text) {
		return WM_NOMATCH
	}
	return WM_MATCH
}

func WildMatch(pattern, text string, flags int) int {
	return doWild([]rune(pattern), []rune(text), flags)
}
