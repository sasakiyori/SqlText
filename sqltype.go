package sqltext

type SqlType int

const (
	Mysql SqlType = iota
	Postgresql
)
