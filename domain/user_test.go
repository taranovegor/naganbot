package domain

import (
	"fmt"
	"testing"
)

type testcaseUser struct {
	id                            int64
	firstName, lastName, username string
	expected                      string
}

var testcasesUserName = []testcaseUser{
	{0, "", "", "", "0"},
	{0, "F", "", "", "F"},
	{0, "", "L", "", "L"},
	{0, "F", "L", "", "F L"},
	{0, "F", "L", "U", "U"},
}

func TestUser_Name(t *testing.T) {
	for _, tc := range testcasesUserName {
		user := NewUser(tc.id, tc.firstName, tc.lastName, tc.username)
		actual := user.Name()
		if tc.expected != actual {
			t.Errorf("given name %s is not equals to expected %s", actual, tc.expected)
		}
	}
}

func TestUser_Mention(t *testing.T) {
	for _, tc := range testcasesUserName {
		user := NewUser(tc.id, tc.firstName, tc.lastName, tc.username)
		var expected string
		if len(tc.username) > 0 {
			expected = fmt.Sprintf("@%s", tc.expected)
		} else {
			expected = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", tc.id, tc.expected)
		}

		actual := user.Mention()
		if expected != actual {
			t.Errorf("given mention %s is not equals to expected %s", actual, expected)
		}
	}
}
