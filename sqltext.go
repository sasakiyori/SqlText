package sqltext

type SqlText interface {
	// RemoveComments remove all the comments of the sql text
	RemoveComments(sql string) string

	// FormatText Minimize the amount of white space in the text and keep only one line
	FormatText(sql string) string

	// CommandType command type of sql text
	CommandType(sql string) (CmdType, error)

	// ReadonlyCommand check out if the sql is readonly command
	ReadonlyCommand(sql string) bool
}
