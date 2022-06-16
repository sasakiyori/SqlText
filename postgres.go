package sqltext

type PostgresqlText struct {
}

func (t PostgresqlText) RemoveComments(sql string) string {
	level := 0 // slash-star comment nested level
	quote := 0 // single quote num
	index, prev, end := 0, 0, len(sql)
	for index < end {
		if sql[index] == '\'' {
			quote++
		} else if quote%2 != 0 {
			index++
		} else if sql[index] == '/' && sql[index+1] == '*' {
			if level == 0 {
				prev = index
			}
			level++
			index += 2
		} else if level > 0 && sql[index] == '*' && sql[index+1] == '/' {
			level--
			index += 2
			if level == 0 {
				sql = sql[:prev] + sql[index:]
				end = len(sql)
				index = prev
				prev = 0
			}
		} else if level == 0 && sql[index] == '-' && sql[index+1] == '-' {
			for i := index + 2; i < end; i++ {
				if sql[i] == '\n' {
					sql = sql[:index] + sql[i:]
					end = len(sql)
					break
				}
			}
		} else {
			index++
		}
	}
	return sql
}

func (t PostgresqlText) FormatText(sql string) string {
	sql = SkipSpacesFromHead(sql)
	quote := 0 // single quote num
	index, end := 0, len(sql)
	for ; index < end; index++ {
		if sql[index] == '\'' {
			quote++
		} else if quote%2 == 0 && IsSpaces(sql[index]) {
			i := index + 1
			for ; i < end; i++ {
				if !IsSpaces(sql[i]) {
					break
				}
			}
			if i == index+1 && sql[index] == ' ' {
				continue
			}
			sql = sql[:index] + " " + sql[i:]
			end = len(sql)
		}
	}
	return sql
}
