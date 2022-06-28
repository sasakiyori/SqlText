package sqltext

import "testing"

func TestNew(t *testing.T) {
	if unknown := New(-9999); unknown != nil {
		t.Error("create a text processor with unknown type should return nil")
	}

	pgsql := New(Postgresql)
	_, ok := pgsql.(*PostgresqlText)
	if !ok {
		t.Error("fail to create a postgresql text processor")
	}
}
