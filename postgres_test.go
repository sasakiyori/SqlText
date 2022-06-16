package sqltext

import (
	"testing"
)

func TestPostgresqlText(t *testing.T) {
	pg := PostgresqlText{}
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "common sql",
			input:  `select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name:   "sql with blank space",
			input:  `   select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name: "sql with crlf",
			input: `
			select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name: "sql with simple comment format",
			input: `-- simple comment format
select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name: "sql with simple slash-star comment",
			input: `/* slash-star comment */
			select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name: "sql with double slash-star comment",
			input: `/* slash-star comment */ /* asd */ -- asdf
select * from test_a;`,
			output: `select * from test_a;`,
		},
		{
			name: "sql with slash-star comment nested comment",
			input: `/* slash-star comment  /* asd */
						*/ select * from test_a;`,
			output: `select * from test_a;`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := pg.FormatText(pg.RemoveComments(tc.input))
			if res != tc.output {
				t.Errorf("expect: %s, actual: %s", tc.output, res)
			}
		})
	}
}
