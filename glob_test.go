package glob_test

import (
	glob "github.com/cirruslabs/go-java-glob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoodCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Pattern        string
		ShouldMatch    []string
		ShouldNotMatch []string
	}{
		// Examples from "Glob with Java NIO" article
		// [1]: https://javapapers.com/java/glob-with-java-nio/
		{"*.txt", []string{"file.txt"}, []string{"file.py"}},
		{"*.{html,htm}", []string{"new.html", "old.htm"}, []string{"file.py"}},
		{"?.txt", []string{"1.txt"}, []string{"22.txt"}},
		{"*.*", []string{"README.md"}, []string{"vmlinuz"}},
		{"/etc/*", []string{"/etc/passwd"}, []string{"/boot/vmlinuz", "/etc/ssh/sshd_config"}},
		{"/home/**", []string{"/home/user1", "/home/user2/www"}, []string{"/"}},
		{"**.go", []string{"dir/file.go", "dir/subdir/file.go"}, []string{"dir/file.py"}},
		{"[ab].txt", []string{"a.txt"}, []string{"c.txt"}},
		{"[a-c].txt", []string{"b.txt"}, []string{"d.txt"}},
		{"[!a].txt", []string{"b.txt"}, []string{"a.txt"}},

		// Potential regression when group matching is re-implemented without checking for path separator
		{"[!a]", []string{"b"}, []string{"/"}},

		// Test cases from OpenJDK's Basic.java[1]
		// [1]: https://github.com/openjdk/jdk/blob/c8de943c1fc3491f6be92ad6b6d959050bbddd44/test/jdk/java/nio/file/PathMatcher/Basic.java#L86-L204
		//
		// Basic
		{"foo.html", []string{"foo.html"}, []string{}},
		{"foo.html", []string{}, []string{"foo.htm"}},
		{"foo.html", []string{}, []string{"bar.html"}},
		//
		// Match zero or more characters
		{"f*", []string{"foo.html"}, []string{}},
		{"*.html", []string{"foo.html"}, []string{}},
		{"foo.html*", []string{"foo.html"}, []string{}},
		{"*foo.html", []string{"foo.html"}, []string{}},
		{"*foo.html*", []string{"foo.html"}, []string{}},
		{"*.htm", []string{}, []string{"foo.html"}},
		{"f.*", []string{}, []string{"foo.html"}},
		//
		// Match one character
		{"?oo.html", []string{"foo.html"}, []string{}},
		{"??o.html", []string{"foo.html"}, []string{}},
		{"???.html", []string{"foo.html"}, []string{}},
		{"???.htm?", []string{"foo.html"}, []string{}},
		{"foo.???", []string{}, []string{"foo.html"}},
		//
		// Group of subpatterns
		{"foo{.html,.class}", []string{"foo.html"}, []string{}},
		{"foo.{class,html}", []string{"foo.html"}, []string{}},
		{"foo{.htm,.class}", []string{}, []string{"foo.html"}},
		//
		// Bracket expressions
		{"[f]oo.html", []string{"foo.html"}, []string{}},
		{"[e-g]oo.html", []string{"foo.html"}, []string{}},
		{"[abcde-g]oo.html", []string{"foo.html"}, []string{}},
		{"[abcdefx-z]oo.html", []string{"foo.html"}, []string{}},
		{"[!a]oo.html", []string{"foo.html"}, []string{}},
		{"[!a-e]oo.html", []string{"foo.html"}, []string{}},
		{"foo[-a-z]bar", []string{"foo-bar"}, []string{}},
		{"foo[!-]html", []string{"foo.html"}, []string{}},
		//
		// Groups of subpattern with bracket expressions
		{"[f]oo.{[h]tml,class}", []string{"foo.html"}, []string{}},
		{"foo.{[a-z]tml,class}", []string{"foo.html"}, []string{}},
		{"foo.{[!a-e]tml,.class}", []string{"foo.html"}, []string{}},
		//
		// Assume special characters are allowed in file names
		{"\\{foo*", []string{"{foo}.html"}, []string{}},
		{"*\\}.html", []string{"{foo}.html"}, []string{}},
		{"\\[foo*", []string{"[foo].html"}, []string{}},
		{"*\\].html", []string{"[foo].html"}, []string{}},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Pattern, func(t *testing.T) {
			t.Parallel()

			re, err := glob.ToRegexPattern(testCase.Pattern, false)
			if err != nil {
				t.Fatal(err)
			}

			for _, shouldMatch := range testCase.ShouldMatch {
				assert.True(t, re.MatchString(shouldMatch))
			}

			for _, shouldMatch := range testCase.ShouldNotMatch {
				assert.False(t, re.MatchString(shouldMatch))
			}
		})
	}
}

func TestBadCases(t *testing.T) {
	testCases := []struct {
		Pattern string
		Error   string
	}{
		// Test cases from OpenJDK's Basic.java[1]
		// [1]: https://github.com/openjdk/jdk/blob/c8de943c1fc3491f6be92ad6b6d959050bbddd44/test/jdk/java/nio/file/PathMatcher/Basic.java#L86-L204
		//
		// Errors
		{"*[a--z]", "invalid range"},
		{"*[a--]", "invalid range"},
		{"*[a-z", "missing ']"},
		{"*{class,java", "missing '}"},
		{"*.{class,{.java}}", "cannot nest groups"},
		{"*.html\\", "no character to escape"},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Pattern, func(t *testing.T) {
			t.Parallel()

			_, err := glob.ToRegexPattern(testCase.Pattern, false)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), testCase.Error)
		})
	}
}
