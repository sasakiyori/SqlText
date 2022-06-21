package sqltext

import (
	"errors"

	pgquery "github.com/pganalyze/pg_query_go/v2"
)

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

func (t PostgresqlText) CommandType(sql string) (CmdType, error) {
	pr, err := pgquery.Parse(sql)
	if err != nil {
		return NIL, err
	}
	stmts := pr.GetStmts()
	if len(stmts) == 0 {
		return NIL, errors.New("invalid sql text")
	}
	node := stmts[0].GetStmt().GetNode()
	// TODO: complement command types
	switch node.(type) {
	case *pgquery.Node_DoStmt:
		return DO, nil
	case *pgquery.Node_LockStmt:
		return LOCK, nil
	case *pgquery.Node_CallStmt:
		return CALL, nil
	case *pgquery.Node_CopyStmt:
		return COPY, nil
	case *pgquery.Node_DropStmt, *pgquery.Node_DropdbStmt, *pgquery.Node_DropRoleStmt,
		*pgquery.Node_DropUserMappingStmt, *pgquery.Node_DropTableSpaceStmt,
		*pgquery.Node_DropSubscriptionStmt, *pgquery.Node_DropOwnedStmt:
		return DROP, nil
	case *pgquery.Node_RuleStmt:
		return RULE, nil
	case *pgquery.Node_ViewStmt:
		return VIEW, nil
	case *pgquery.Node_AlterDomainStmt, *pgquery.Node_AlterCollationStmt, *pgquery.Node_AlterDatabaseSetStmt,
		*pgquery.Node_AlterDatabaseStmt, *pgquery.Node_AlterDefaultPrivilegesStmt, *pgquery.Node_AlterEnumStmt,
		*pgquery.Node_AlterEventTrigStmt, *pgquery.Node_AlterExtensionContentsStmt, *pgquery.Node_AlterExtensionStmt,
		*pgquery.Node_AlterFdwStmt, *pgquery.Node_AlterForeignServerStmt, *pgquery.Node_AlterFunctionStmt,
		*pgquery.Node_AlterObjectDependsStmt, *pgquery.Node_AlterObjectSchemaStmt, *pgquery.Node_AlterOpFamilyStmt,
		*pgquery.Node_AlterOperatorStmt, *pgquery.Node_AlterOwnerStmt, *pgquery.Node_AlterPolicyStmt,
		*pgquery.Node_AlterPublicationStmt, *pgquery.Node_AlterRoleSetStmt, *pgquery.Node_AlterRoleStmt,
		*pgquery.Node_AlterSeqStmt, *pgquery.Node_AlterStatsStmt, *pgquery.Node_AlterSubscriptionStmt,
		*pgquery.Node_AlterSystemStmt, *pgquery.Node_AlterTableCmd, *pgquery.Node_AlterTableMoveAllStmt,
		*pgquery.Node_AlterTableSpaceOptionsStmt, *pgquery.Node_AlterTableStmt, *pgquery.Node_AlterTsconfigurationStmt,
		*pgquery.Node_AlterTsdictionaryStmt, *pgquery.Node_AlterTypeStmt, *pgquery.Node_AlterUserMappingStmt,
		*pgquery.Node_AlternativeSubPlan:
		return ALTER, nil
	case *pgquery.Node_FetchStmt:
		return FETCH, nil
	case *pgquery.Node_IndexStmt, *pgquery.Node_ReindexStmt:
		return INDEX, nil
	case *pgquery.Node_GrantStmt, *pgquery.Node_GrantRoleStmt:
		return GRANT, nil
	case *pgquery.Node_CreateStmt, *pgquery.Node_CreatedbStmt, *pgquery.Node_CreateAmStmt,
		*pgquery.Node_CreateCastStmt, *pgquery.Node_CreateConversionStmt, *pgquery.Node_CreateDomainStmt,
		*pgquery.Node_CreateEnumStmt, *pgquery.Node_CreateEventTrigStmt, *pgquery.Node_CreateExtensionStmt,
		*pgquery.Node_CreateFdwStmt, *pgquery.Node_CreateForeignServerStmt, *pgquery.Node_CreateForeignTableStmt,
		*pgquery.Node_CreateFunctionStmt, *pgquery.Node_CreateOpClassItem, *pgquery.Node_CreateOpClassStmt,
		*pgquery.Node_CreateOpFamilyStmt, *pgquery.Node_CreatePlangStmt, *pgquery.Node_CreatePolicyStmt,
		*pgquery.Node_CreatePublicationStmt, *pgquery.Node_CreateRangeStmt, *pgquery.Node_CreateRoleStmt,
		*pgquery.Node_CreateSchemaStmt, *pgquery.Node_CreateSeqStmt, *pgquery.Node_CreateStatsStmt,
		*pgquery.Node_CreateSubscriptionStmt, *pgquery.Node_CreateTableAsStmt, *pgquery.Node_CreateTableSpaceStmt,
		*pgquery.Node_CreateTransformStmt, *pgquery.Node_CreateTrigStmt, *pgquery.Node_CreateUserMappingStmt:
		return CREATE, nil
	case *pgquery.Node_SelectStmt:
		return SELECT, nil
	case *pgquery.Node_UpdateStmt:
		return UPDATE, nil
	case *pgquery.Node_InsertStmt:
		return INSERT, nil
	case *pgquery.Node_DeleteStmt:
		return DELETE, nil
	case *pgquery.Node_DeclareCursorStmt:
		return CURSOR, nil
	case *pgquery.Node_ExplainStmt:
		return EXPLAIN, nil
	case *pgquery.Node_PrepareStmt:
		return PREPARE, nil
	case *pgquery.Node_ExecuteStmt:
		return EXECUTE, nil
	case *pgquery.Node_TruncateStmt:
		return TRUNCATE, nil
	case *pgquery.Node_CheckPointStmt:
		return CHECKPOINT, nil
	case *pgquery.Node_TransactionStmt:
		return TRANSACTION, nil
	default:
		return NIL, errors.New("unknown sql command type")
	}
}

func (t PostgresqlText) ReadonlyCommand(sql string) bool {
	ct, err := t.CommandType(sql)
	if err != nil {
		return false
	}
	// TODO: EXPLAIN -- analyze or not
	// TODO: CALL    -- unknown procedure to be invoked
	// DOUBT: fetch changes the position of cursor, is it can be called a readonly command?
	switch ct {
	case FETCH, SELECT, EXPLAIN, CHECKPOINT:
		return true
	default:
		return false
	}
}
