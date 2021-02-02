// Copyright (c) 2008, 2009, Oracle and/or its affiliates. All rights reserved.
// DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
//
// This code is free software; you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 2 only, as
// published by the Free Software Foundation.  Oracle designates this
// particular file as subject to the "Classpath" exception as provided
// by Oracle in the LICENSE file that accompanied this code.
//
// This code is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
// version 2 for more details (a copy is included in the LICENSE file that
// accompanied this code).
//
// You should have received a copy of the GNU General Public License version
// 2 along with this work; if not, write to the Free Software Foundation,
// Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301 USA.
//
// Please contact Oracle, 500 Oracle Parkway, Redwood Shores, CA 94065 USA
// or visit www.oracle.com if you need additional information or have any
// questions.

package glob

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const regexMetaChars = ".^$+{[]|()"
const globMetaChars = "\\*?[{"

const EOL = 0

var ErrPatternSyntax = errors.New("invalid pattern syntax")

func isRegexMeta(c uint8) bool {
	return strings.ContainsRune(regexMetaChars, rune(c))
}

func isGlobMeta(c uint8) bool {
	return strings.ContainsRune(globMetaChars, rune(c))
}

func next(glob string, i int) uint8 {
	if i < len(glob) {
		return glob[i]
	}

	return EOL
}

// nolint:gocognit,gocyclo // keeping the structure as is makes it easier to spot the differences with Glob.java
func ToRegexPattern(globPattern string, isDos bool) (*regexp.Regexp, error) {
	var inGroup bool
	regex := "^"

	var i int

	for i < len(globPattern) {
		c := globPattern[i]; i++
		switch c {
		case '\\':
			// escape special characters
			if i == len(globPattern) {
				return nil, fmt.Errorf("%w: no character to escape in %s at %d", ErrPatternSyntax, globPattern, i-1)
			}
			next := globPattern[i]; i++
			if isGlobMeta(next) || isRegexMeta(next) {
				regex += "\\"
			}
			regex += string(next)
		case '/':
			if isDos {
				regex += "\\\\"
			} else {
				regex += string(c)
			}
		case '[':
			// don't match name separator in class
			if isDos {
				regex += "[[^\\\\]&&["
			} else {
				regex += "[[^/]&&["
			}
			if next(globPattern, i) == '^' {
				// escape the regex negation char if it appears
				regex += "\\^"
				i++
			} else {
				// negation
				if next(globPattern, i) == '!' {
					regex += "^"
					i++
				}
				// hyphen allowed at start
				if next(globPattern, i) == '-' {
					regex += "-"
					i++
				}
			}

			var hasRangeStart bool
			var last uint8

			for i < len(globPattern) {
				c = globPattern[i]; i++
				if c == ']' {
					break
				}
				if c == '/' || (isDos && c == '\\') {
					return nil, fmt.Errorf("%w: explicit 'name separator' in class in %s at %d", ErrPatternSyntax, globPattern, i-1)
				}
				// TBD: how to specify ']' in a class?
				if c == '\\' || c == '[' || c == '&' && next(globPattern, i) == '&' {
					// escape '\', '[' or "&&" for regex class
					regex += "\\"
				}
				regex += string(c)

				if c == '-' {
					if !hasRangeStart {
						return nil, fmt.Errorf("%w: invalid range in %s at %d", ErrPatternSyntax, globPattern, i-1)
					}
					c = next(globPattern, i); i++
					if c == EOL || c == ']' {
						break
					}
					if c < last {
						return nil, fmt.Errorf("%w: invalid range in %s at %d", ErrPatternSyntax, globPattern, i-3)
					}
					regex += string(c)
					hasRangeStart = false
				} else {
					hasRangeStart = true
					last = c
				}
			}
			if c != ']' {
				return nil, fmt.Errorf("%w: missing '] in %s at %d", ErrPatternSyntax, globPattern, i-1)
			}
			regex += "]]"
		case '{':
			if inGroup {
				return nil, fmt.Errorf("%w: cannot nest groups in %s at %d", ErrPatternSyntax, globPattern, i-1)
			}
			regex += "(?:(?:"
			inGroup = true
		case '}':
			if inGroup {
				regex += "))"
				inGroup = false
			} else {
				regex += "}"
			}
		case ',':
			if inGroup {
				regex += ")|(?:"
			} else {
				regex += ","
			}
		case '*':
			if next(globPattern, i) == '*' {
				// crosses directory boundaries
				regex += ".*"
				i++
			} else {
				// within directory boundary
				if isDos {
					regex += "[^\\\\]*"
				} else {
					regex += "[^/]*"
				}
			}
		case '?':
			if isDos {
				regex += "[^\\\\]"
			} else {
				regex += "[^/]"
			}
		default:
			if isRegexMeta(c) {
				regex += "\\"
			}
			regex += string(c)
		}
	}

	if inGroup {
		return nil, fmt.Errorf("%w: missing '} in %s at %d", ErrPatternSyntax, globPattern, i-1)
	}

	return regexp.Compile(regex + "$")
}
