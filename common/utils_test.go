package common

import "testing"

func TestSublimeContains(t *testing.T) {
	tests := []struct {
		text   string
		substr string
		pass   bool
	}{
		{"hello", "lo", true},
		{"abcdefg", "cf", true},
		{"abcdefg", "a", true},
		{"abcdefg", "b", true},
		{"abcdefg", "cfa", false},
		{"abcdefg", "aa", false},
		{"世界", "a", false},
		{"Hello 世界", "界", true},
		{"Hello 世界", "elo", true},
	}
	for _, v := range tests {
		res := SublimeContains(v.text, v.substr)
		if res != v.pass {
			t.Fatalf("Failed: %v - res:%v", v, res)
		}
	}
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		orig   string
		expect string
	}{
		// {"C:\\hello", "C:/hello"}, // Only works in windows
		{"", "."},
		{"//../foo", "/foo"},
		{"/../../", "/"},
		{"/hello/world/..", "/hello"},
		{"/..", "/"},
		{"/foo/..", "/"},
		{"/-/foo", "/-/foo"},
	}
	for _, v := range tests {
		res := CleanPath(v.orig)
		if res != v.expect {
			t.Fatalf("Clean path(%v) expect(%v) but got(%v)", v.orig, v.expect, res)
		}
	}
}

func TestCheckPath(t *testing.T) {
	tests := []struct {
		index int
		path1 string
		path2 string
		want  bool
	}{
		{1, "a/b/c/*/**/*/d/e/f", "a/b/c/d/e/f", false},
		{2, "a/b/c/*/**/*/d/e/f", "a/b/c/t/d/e/f", true},
		{3, "var/users/dingyubo@dp.tech/application_data/*/**/*", "var/users/dingyubo@dp.tech/application_data", true},
		{4, "var/users/dingyubo@dp.tech/application_data/*/**/*", "var/users/dingyubo@dp.tech/application_data/a", true},
		{5, "var/users/dingyubo@dp.tech/application_data/*/**/*", "var/users/dingyubo@dp.tech/application_data/a/b", true},
		{6, "var/users/dingyubo@dp.tech/application_data/*/**/*", "var/users/dingyubo@dp.tech/application_data/a/b/c/d/e/f/g", true},
		{7, "var/users/dingyubo@dp.tech/application_data/*/**/*/a", "var/users/dingyubo@dp.tech/application_data/a/a", true},
		{8, "a/b/c/", "/a/b/c/", true},
		{9, "a/*/b", "a/b/b", true},
		{10, "a/*/b/c", "a/b/b", false},
		{11, "a/*/b/*", "a/b/b", true},
		{12, "a/*/b/**", "a/b/b", true},
		{13, "a/*/b/**", "a/b/b/", true},
		{14, "**", "", true},
		{15, "**/*", "a/b", true},
		{16, "**/*", "a/b/c/d/e/f", true},
		{17, "**", "a/b/c", true},
		{18, "*/a/*/b/*/c", "a/b/c", false},
		{19, "*/a/*/b/*/c", "a/a/b/c", false},
		{20, "*/a/*/b/**/c", "a/a/b/b/c", true},
		{21, "*/a/*/b/*/c", "a/a/b/b/c", false},
	}

	for _, tt := range tests {
		if got := CheckPath(tt.path1, tt.path2); got != tt.want {
			t.Errorf("checkPath(%d): checkPath(%q, %q) = %v; want %v", tt.index, tt.path1, tt.path2, got, tt.want)
		}
	}
}
