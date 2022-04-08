package wildmatch

import "testing"

func Test_WildMatch(t *testing.T) {
	type args struct {
		text    string
		pattern string
		flags   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"match normal",
			args{
				"foo",
				"foo",
				0,
			},
			0,
		},
		{
			"match fail",
			args{
				"foo",
				"bar",
				0,
			},
			1,
		},
		{
			"match nil",
			args{
				"",
				"",
				0,
			},
			0,
		},
		{
			"match noraml ?",
			args{
				"foo",
				"???",
				0,
			},
			0,
		},
		{
			"match fail ?",
			args{
				"foo",
				"??",
				0,
			},
			1,
		},
		{
			"match noraml *",
			args{
				"foo",
				"*",
				0,
			},
			0,
		},
		{
			"match suffix *",
			args{
				"foo",
				"f*",
				0,
			},
			0,
		},
		{
			"match suffix extra *",
			args{
				"foo",
				"foo*",
				0,
			},
			0,
		},
		{
			"match prefix *",
			args{
				"aaaaaaabababab",
				"*ab",
				0,
			},
			0,
		},
		{
			"match prefix * fail",
			args{
				"foo",
				"*f",
				0,
			},
			1,
		},
		{
			"match surround *",
			args{
				"foo",
				"*foo*",
				0,
			},
			0,
		},
		{
			"match interval *",
			args{
				"foobar",
				"*ob*a*r*",
				0,
			},
			0,
		},
		{
			"match escape *",
			args{
				`foo*`,
				`foo\*`,
				0,
			},
			0,
		},
		{
			"match escape * fail",
			args{
				`foobar`,
				`foo\*bar`,
				0,
			},
			1,
		},
		{
			"match escape",
			args{
				`f\oo`,
				`f\\oo`,
				0,
			},
			0,
		},
		{
			"match brackets",
			args{
				`ten`,
				`[ten]`,
				0,
			},
			1,
		},
		{
			"match brackets fail",
			args{
				`ten`,
				`[ten]`,
				0,
			},
			1,
		},
		{
			"match brackets ?",
			args{
				`ball`,
				`*[al]?`,
				0,
			},
			0,
		},
		{
			"match brackets invert",
			args{
				`ten`,
				`**[!te]`,
				0,
			},
			0,
		},
		{
			"match brackets invert fail",
			args{
				`ten`,
				`**[!ten]`,
				0,
			},
			-1,
		},
		{
			"match brackets hyphen",
			args{
				`ten`,
				`t[a-g]n`,
				0,
			},
			0,
		},
		{
			"match brackets hyphen invert !",
			args{
				`ton`,
				`t[!a-g]n`,
				0,
			},
			0,
		},
		{
			"match brackets hyphen invert ! fail",
			args{
				`ten`,
				`t[!a-g]n`,
				0,
			},
			1,
		},
		{
			"match brackets hyphen invert ^",
			args{
				`ton`,
				`t[^a-g]n`,
				0,
			},
			0,
		},
		{
			"match brackets partial bracket",
			args{
				`a]b`,
				`a[]]b`,
				0,
			},
			0,
		},
		{
			"match brackets partial bracket hyphen",
			args{
				`a-b`,
				`a[]-]b`,
				0,
			},
			0,
		},
		{
			"match brackets partial bracket hyphen 2",
			args{
				`a]b`,
				`a[]-]b`,
				0,
			},
			0,
		},
		{
			"match brackets partial bracket hyphen fail",
			args{
				`aab`,
				`a[]-]b`,
				0,
			},
			1,
		},
		{
			"match brackets partial bracket hyphen 3",
			args{
				`aab`,
				`a[]a-]b`,
				0,
			},
			0,
		},
		{
			"match brackets partial bracket only",
			args{
				`]`,
				`]`,
				0,
			},
			0,
		},
		// Extended slash-matching features
		{
			"match slash asterisk",
			args{
				`foo/baz/bar`,
				`foo*bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk",
			args{
				`foo/baz/bar`,
				`foo**bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 2",
			args{
				`foobazbar`,
				`foo**bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/baz/bar`,
				`foo/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/baz/bar`,
				`foo/**/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/b/a/z/bar`,
				`foo/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/b/a/z/bar`,
				`foo/**/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`foo/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`foo/**/**/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`foo?bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`foo[/]bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`foo[^a-z]bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar`,
				`f[^eiu][^eiu][^eiu][^eiu][^eiu]r`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo-bar`,
				`f[^eiu][^eiu][^eiu][^eiu][^eiu]r`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo`,
				`**/foo`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`XXX/foo`,
				`**/foo`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`bar/baz/foo`,
				`**/foo`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`bar/baz/foo`,
				`*/foo`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar/baz`,
				`**/bar*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/baz`,
				`**/bar/*`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/baz/`,
				`**/bar/*`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/baz/`,
				`**/bar/**`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar`,
				`**/bar/*`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/`,
				`**/bar/**`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar/baz`,
				`**/bar**`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`foo/bar/baz/x`,
				`**/bar/**`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/baz/x`,
				`*/bar/**`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match slash double asterisk 3",
			args{
				`deep/foo/bar/baz/x`,
				`**/bar/*/*`,
				WM_PATHNAME,
			},
			0,
		},
		// Character class tests
		{
			"match Character class alpha digit upper",
			args{
				`a1B`,
				`[[:alpha:]][[:digit:]][[:upper:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class digit upper space",
			args{
				`a`,
				`[[:digit:][:upper:][:space:]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match Character class digit upper space",
			args{
				`A`,
				`[[:digit:][:upper:][:space:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class digit upper space",
			args{
				`1`,
				`[[:digit:][:upper:][:space:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class digit upper space",
			args{
				`1`,
				`[[:digit:][:upper:][:spaci:]]`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match Character class digit upper space",
			args{
				` `,
				`[[:digit:][:upper:][:space:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class digit upper space",
			args{
				`.`,
				`[[:digit:][:upper:][:space:]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match Character class digit upper space",
			args{
				`.`,
				`[[:digit:][:punct:][:space:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`5`,
				`[[:xdigit:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`f`,
				`[[:xdigit:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`D`,
				`[[:xdigit:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`_`,
				`[[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:graph:][:lower:][:print:][:punct:][:space:][:upper:][:xdigit:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`.`,
				`[^[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:lower:][:space:][:upper:][:xdigit:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`5`,
				`[a-c[:digit:]x-z]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`b`,
				`[a-c[:digit:]x-z]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`y`,
				`[a-c[:digit:]x-z]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match Character class xdigit",
			args{
				`q`,
				`[a-c[:digit:]x-z]`,
				WM_PATHNAME,
			},
			1,
		},
		// Additional tests, including some malformed wildmatch patterns
		{
			"match additional",
			args{
				`]`,
				`[\\-^]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`[`,
				`[\\-^]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`-`,
				`[\-_]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`]`,
				`[\]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`\]`,
				`[\]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`\`,
				`[\]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`ab`,
				`a[]b`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`a[]b`,
				`a[]b`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`ab[`,
				`ab[`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`ab`,
				`[!`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`ab`,
				`[-`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`-`,
				`[-]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-`,
				`[a-`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`-`,
				`[!a-`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match additional",
			args{
				`-`,
				`[--A]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`5`,
				`[--A]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				` `,
				`[ --]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`$`,
				`[ --]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-`,
				`[ --]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`0`,
				`[ --]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`-`,
				`[---]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-`,
				`[-----]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`j`,
				`[a-e-n]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`-`,
				`[a-e-n]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`a`,
				`[!------]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`[`,
				`[]-a]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`^`,
				`[]-a]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`^`,
				`[!]-a]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`[`,
				`[!]-a]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`^`,
				`[a^bc]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-b]`,
				`[a-]b]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`\`,
				`[\]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`\`,
				`[\\]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`\`,
				`[!\\]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`G`,
				`[A-\\]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`aaabbb`,
				`b*a`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`aabcaa`,
				`*ba*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`,`,
				`[,]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`,`,
				`[\\,]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`\`,
				`[\\,]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-`,
				`[,-.]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`+`,
				`[,-.]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`-.]`,
				`[,-.]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`2`,
				`[\1-\3]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`3`,
				`[\1-\3]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`4`,
				`[\1-\3]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match additional",
			args{
				`\`,
				`[[-\]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`[`,
				`[[-\]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`]`,
				`[[-\]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match additional",
			args{
				`-`,
				`[[-\]]`,
				WM_PATHNAME,
			},
			1,
		},
		// Test recursion
		{
			"match recursion",
			args{
				`-adobe-courier-bold-o-normal--12-120-75-75-m-70-iso8859-1`,
				`-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`-adobe-courier-bold-o-normal--12-120-75-75-X-70-iso8859-1`,
				`-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
				WM_PATHNAME,
			},
			-1,
		},
		{
			"match recursion",
			args{
				`-adobe-courier-bold-o-normal--12-120-75-75-/-70-iso8859-1`,
				`-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`XXX/adobe/courier/bold/o/normal//12/120/75/75/m/70/iso8859/1`,
				`XXX/*/*/*/*/*/*/12/*/*/*/m/*/*/*`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`XXX/adobe/courier/bold/o/normal//12/120/75/75/X/70/iso8859/1`,
				`XXX/*/*/*/*/*/*/12/*/*/*/m/*/*/*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txt`,
				`**/*a*b*g*n*t`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txtz`,
				`**/*a*b*g*n*t`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`foo`,
				`*/*/*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`foo/bar`,
				`*/*/*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`foo/bba/arr`,
				`*/*/*`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`foo/bb/aa/rr`,
				`*/*/*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`foo/bb/aa/rr`,
				`**/**/**`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`abcXdefXghi`,
				`*X*i`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`ab/cXd/efXg/hi`,
				`*X*i`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match recursion",
			args{
				`ab/cXd/efXg/hi`,
				`*/*X*/*/*i`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match recursion",
			args{
				`ab/cXd/efXg/hi`,
				`**/*X*/**/*i`,
				WM_PATHNAME,
			},
			0,
		},
		// Extra pathmatch tests
		{
			"match pathmatch",
			args{
				`foo`,
				`fo`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bar`,
				`foo/bar`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match pathmatch",
			args{
				`foo/bar`,
				`foo/*`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/**`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo*`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo**`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/*arr`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/**arr`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/*z`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bba/arr`,
				`foo/**z`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bar`,
				`foo?bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bar`,
				`foo[/]bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`foo/bar`,
				`foo[^a-z]bar`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match pathmatch",
			args{
				`ab/cXd/efXg/hi`,
				`*Xg*i`,
				WM_PATHNAME,
			},
			1,
		},
		// Extra case-sensitivity tests
		{
			"match case-sensitivity",
			args{
				`a`,
				`[A-Z]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[A-Z]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[a-z]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`a`,
				`[a-z]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`a`,
				`[[:upper:]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[[:upper:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[[:lower:]]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`a`,
				`[[:lower:]]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[B-Za]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`a`,
				`[B-Za]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`A`,
				`[B-a]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`a`,
				`[B-a]`,
				WM_PATHNAME,
			},
			0,
		},
		{
			"match case-sensitivity",
			args{
				`z`,
				`[Z-y]`,
				WM_PATHNAME,
			},
			1,
		},
		{
			"match case-sensitivity",
			args{
				`Z`,
				`[Z-y]`,
				WM_PATHNAME,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildMatch(tt.args.pattern, tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("doWild() = %v, want %v", got, tt.want)
			}
		})
	}
}
