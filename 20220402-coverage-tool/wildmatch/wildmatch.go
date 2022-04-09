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
func isBlank(c rune) bool {
	return c == ' ' || (c) == '\t'
}

func isAlphaNum(c rune) bool {
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
			tCh, pCh, prevCh rune
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
			if tCh != pCh {
				return WM_NOMATCH
			}
		case '?':
			/* Match anything but '/'. */
			if (flags&WM_PATHNAME) != 0 && tCh == '/' {
				return WM_NOMATCH
			}
		case '*':
			iPattern++

			matched := 0
			matchSlash := false

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
					matchSlash = true
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

					matchSlash = true
				} else { /* WM_PATHNAME is set */
					matchSlash = false
				}
			} else {
				/* without WM_PATHNAME, '*' == '**' */
				if (flags & WM_PATHNAME) != 0 {
					matchSlash = false
				} else {
					matchSlash = true
				}
			}

			if iPattern == len(pattern) {
				/* Trailing "**" matches everything.  Trailing "*" matches
				 * only if there are no more slash characters. */
				if !matchSlash {
					if strchr(text[iText:], '/') != len(text[iText:]) {
						return WM_NOMATCH
					}
				}
				return WM_MATCH
			} else if !matchSlash && pattern[iPattern] == '/' {
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
						if !matchSlash && tCh == '/' {
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
					if matchSlash || matched != WM_ABORT_TO_STARSTAR {
						return matched
					}
				} else if !matchSlash && tCh == '/' {
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

			matched, negated := false, false

			if iPattern < len(pattern) {
				pCh = pattern[iPattern]
			}
			if pCh == NEGATE_CLASS2 {
				pCh = NEGATE_CLASS
			}

			/* Assign literal 1/0 because of "matched" comparison. */
			if pCh == NEGATE_CLASS {
				negated = true
			}

			if negated { /* Inverted character class. */
				iPattern++
				if iPattern < len(pattern) {
					pCh = pattern[iPattern]
				}
			}

			prevCh = 0
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

					matched = (tCh == pCh)
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
						matched = true
					} else if (flags&WM_CASEFOLD) != 0 && unicode.IsLower(tCh) {
						tChUpper := unicode.ToUpper(tCh)
						if tChUpper <= pCh && tChUpper >= prevCh {
							matched = true
						}
					}
					pCh = 0 /* This makes "prev_ch" get set to 0. */
				} else if pCh == '[' && pattern[iPattern+1] == ':' {
					chBeg := iPattern + 2
					for {
						pCh = pattern[iPattern]
						if pCh == ']' || iPattern == len(pattern) {
							break
						}
						iPattern++
					}
					if iPattern == len(pattern) {
						return WM_ABORT_ALL
					}

					chLen := iPattern - chBeg - 1
					if chLen < 0 || pattern[iPattern-1] != ':' {
						pattern = pattern[chBeg-2:]
						pCh = '['
						matched = (tCh == pCh)

						continue
					}

					ch := pattern[chBeg : chBeg+chLen]
					if equals(ch, []rune("alnum")) {
						if isAlphaNum(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("alpha")) {
						if unicode.IsLetter(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("blank")) {
						if isBlank(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("cntrl")) {
						if unicode.IsControl(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("digit")) {
						if unicode.IsDigit(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("graph")) {
						if unicode.IsGraphic(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("lower")) {
						if unicode.IsLower(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("print")) {
						if unicode.IsPrint(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("punct")) {
						if unicode.IsPunct(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("space")) {
						if unicode.IsSpace(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("upper")) {
						if unicode.IsUpper(tCh) {
							matched = true
						} else if (flags&WM_CASEFOLD) != 0 && unicode.IsLower(tCh) {
							matched = true
						}
					} else if equals(ch, []rune("xdigit")) {
						if unicode.Is(unicode.Hex_Digit, tCh) {
							matched = true
						}
					} else { /* malformed [:class:] string */
						return WM_ABORT_ALL
					}
					pCh = 0 /* This makes "prev_ch" get set to 0. */
				} else if tCh == pCh {
					matched = true
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
		default:
			if tCh != pCh {
				return WM_NOMATCH
			}
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
