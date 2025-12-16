package router

import (
	"regexp"
	"testing"
)

func TestRegexps(t *testing.T) {
	tests := map[string]struct {
		reg       string
		match     []string
		dontMatch []string
	}{
		"digit": {
			reg: DigitRegExp,
			match: []string{
				"0",
				"123",
				"999999",
			},
			dontMatch: []string{
				"",
				"abc",
				"12a",
				"-123",
			},
		},

		"uuid v4": {
			reg: UuidV4Regexp,
			match: []string{
				"550e8400-e29b-41d4-a716-446655440000",
				"f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			dontMatch: []string{
				"550e8400-e29b-11d4-a716-446655440000",
				"not-a-uuid",
				"550e8400e29b41d4a716446655440000",
			},
		},

		"email": {
			reg: EmailRegexp,
			match: []string{
				"test@example.com",
				"user+tag@domain.co",
			},
			dontMatch: []string{
				"test@",
				"@example.com",
				"test.example.com",
			},
		},

		"string": {
			reg: StringRegExp,
			match: []string{
				"abc",
				"hello-world",
				"with space",
			},
			dontMatch: []string{
				"/",
				"abc/def",
			},
		},

		"token": {
			reg: TokenRegexp,
			match: []string{
				"abc",
				"ABC123",
				"token42",
			},
			dontMatch: []string{
				"with-dash",
				"with space",
				"",
			},
		},

		"multi token": {
			reg: MultiTokenRegexp,
			match: []string{
				"one",
				"one;two",
				"one_two;three4",
			},
			dontMatch: []string{
				";",
				"one;",
				"one;;two",
			},
		},

		"slug": {
			reg: SlugRegexp,
			match: []string{
				"simple-slug",
				"with_underscore",
				"v1.2.3",
			},
			dontMatch: []string{
				"with space",
				"/slash",
			},
		},

		"jwt": {
			reg: JwtRegexp,
			match: []string{
				"aaa.bbb.ccc",
				"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.signature",
			},
			dontMatch: []string{
				"aaa.bbb",
				"aaa.bbb.ccc.ddd",
			},
		},

		"base64": {
			reg: Base64Regexp,
			match: []string{
				"",
				"Zg==",
				"Zm9v",
				"Zm9vYmFy",
			},
			dontMatch: []string{
				"not base64",
				"Z===",
			},
		},

		"path": {
			reg: PathRegexp,
			match: []string{
				"file.txt",
				"dir/subdir/file",
				"dir.with.dots/file_name",
				"Archive.zip",
			},
			dontMatch: []string{
				"/absolute/path",
				"dir//file",
				"dir/../file",
			},
		},

		"any": {
			reg: AnyRegexp,
			match: []string{
				"",
				"anything",
				"/with/slash",
			},
			dontMatch: []string{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			re, err := regexp.Compile("^" + tt.reg + "$")
			if err != nil {
				t.Fatalf("invalid regexp %q: %v", tt.reg, err)
			}

			for _, s := range tt.match {
				sub := re.FindStringSubmatch(s)
				if sub == nil {
					t.Fatalf("expected %q to match %q", s, tt.reg)
				}

				if sub[0] != s {
					t.Fatalf("full match mismatch: got %q want %q", sub[0], s)
				}

				if len(sub) < 2 {
					t.Fatalf("regexp %q has no capture group", tt.reg)
				}

				if sub[1] != s {
					t.Fatalf(
						"capture mismatch: got %q want %q (regexp %q)",
						sub[1],
						s,
						tt.reg,
					)
				}
			}

			for _, s := range tt.dontMatch {
				if re.MatchString(s) {
					t.Errorf("expected %q NOT to match %q", s, tt.reg)
				}
			}
		})
	}
}
