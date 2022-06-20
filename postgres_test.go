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
				t.Errorf("FormatText expect: %s, actual: %s", tc.output, res)
			}
		})
	}
}

func TestPostgresqlCommandType(t *testing.T) {
	pg := PostgresqlText{}
	testCases := []struct {
		name     string
		sql      string
		cmdType  CmdType
		readonly bool
	}{
		{"do", "DO $$DECLARE r record;$$;", DO, false},
		{"lock", "lock table films in share mode;", LOCK, false},
		{"call", "call dosth();", CALL, false},
		{"copy", "copy test_a from '/path/to/file';", COPY, false},
		{"drop", "drop table test_a;", DROP, false},
		{"rule", "create rule rule_a as on update to test_a do insert into test_a values(1,'x');", RULE, false},
		{"view", "create view vista as select 'hello world';", VIEW, false},
		{"alter", "alter table test_a add column new_col varchar(30);", ALTER, false},
		{"fetch", "fetch forward 5 from liahona;", FETCH, true},
		{"index", "create unique index title_idx on films (title);", INDEX, false},
		{"grant", "grant insert on films to public;", GRANT, false},
		{"create", "create table test_a (id int, name text);", CREATE, false},
		{"select", "select * from test_a;", SELECT, true},
		{"update", "update test_a set name='x' where id=1;", UPDATE, false},
		{"insert", "insert into test_a values (1, 'a')", INSERT, false},
		{"delete", "delete from test_a where id=1;", DELETE, false},
		{"cursor", "declare liahona scroll cursor for select * from films;", CURSOR, false},
		{"explain", "explain select * from test_a;", EXPLAIN, true},
		{"prepare", "prepare foo (int,text) as insert into test_a values($1,$2);", PREPARE, false},
		{"execute", "execute foo(1,'a');", EXECUTE, false},
		{"truncate", "truncate table test_a;", TRUNCATE, false},
		{"checkpoint", "checkpoint;", CHECKPOINT, true},
		{"transaction", "begin;", TRANSACTION, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if ct, err := pg.CommandType(tc.sql); err != nil || ct != tc.cmdType {
				t.Errorf("CommandType expect: %d, actual: %d, error: %v", tc.cmdType, ct, err)
			}
			if ro := pg.ReadonlyCommand(tc.sql); ro != tc.readonly {
				t.Errorf("ReadonlyCommand expect: %t, actual: %t", tc.readonly, ro)
			}
		})
	}
}
