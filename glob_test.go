package glob_test

import (
	"fmt"
	glob "github.com/cirruslabs/go-java-glob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicPatterns(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Pattern     string
		ShouldMatch []string
		ShouldNotMatch []string
	}{
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
		{"[!a]", []string{"b"}, []string{"/"}},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Pattern, func(t *testing.T) {
			re, err := glob.ToRegexPattern(testCase.Pattern, false)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(re.String())

			for _, shouldMatch := range testCase.ShouldMatch {
				assert.True(t, re.MatchString(shouldMatch))
			}

			for _, shouldMatch := range testCase.ShouldNotMatch {
				assert.False(t, re.MatchString(shouldMatch))
			}
		})
	}
}
