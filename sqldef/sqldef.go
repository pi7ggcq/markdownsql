package sqldef

import (
	"fmt"
	"regexp"
	"strings"
)

type SQLDefinition struct {
	DatabaseName string
	Tables       []TableDefinition
}

func (sd SQLDefinition) ToSQLStmt() string {
	sql := ``
	for _, tbl := range sd.Tables {
		// fmt.Printf("デバッグ >> テーブル名 %v\n", tbl.PrimaryKeys)
		sql = sql + fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", tbl.Name)
		sql = sql + fmt.Sprintf("CREATE TABLE `%s` (\n", tbl.Name)

		length := len(tbl.Columns)
		for i, column := range tbl.Columns {
			comma := `,`
			if !tbl.HasIndexes() && i == length-1 {
				comma = ``
			}

			sql = sql + fmt.Sprintf("  %s%s\n", column.ToSQLStmt(), comma)
		}

		comma := `,`
		fkLen := len(tbl.ForeignKeys)
		indexLen := len(tbl.Indexes)
		if (len(tbl.PrimaryKeys) == 0) {
			panic("Primary Key が設定されていません！")
		}
		for _, pk := range tbl.PrimaryKeys {
			prefix := "\n"
			if 0 < (fkLen + indexLen) {
				prefix = fmt.Sprintf("%s\n", comma)
			}
			sql = sql + fmt.Sprintf("  %s%s", pk.ToSQLStmt(), prefix)
		}

		for i, fk := range tbl.ForeignKeys {
			prefix := "\n"
			if indexLen != 0 || i != fkLen - 1 {
				prefix = fmt.Sprintf("%s\n", comma)
			}
			sql = sql + fmt.Sprintf("  %s%s", fk.ToSQLStmt(), prefix)
		}

		length = len(tbl.Indexes)
		for i, index := range tbl.Indexes {
			comma := `,`
			if i == indexLen - 1 {
				comma = ``
			}
			sql = sql + fmt.Sprintf("  %s%s\n", index.ToSQLStmt(), comma)
		}

		sql = sql + fmt.Sprintln(`) ENGINE = InnoDB DEFAULT CHARSET utf8;`)
		sql = sql + fmt.Sprintln(``)
	}

	return sql
}

type TableDefinition struct {
	Name    string
	Columns []ColumnDefinition
	Indexes []IndexDefinition
	PrimaryKeys []PrimaryKeysDefinition
	ForeignKeys []ForeignKeysDefinition
}

func (tbl TableDefinition) HasIndexes() bool {
	return 0 < len(tbl.Indexes)
}

type ColumnDefinition struct {
	Desc	string
	Name    string
	Type    string
	Null    bool
	Default string
	Key     string
	Comment string
}

func (cd ColumnDefinition) ToSQLStmt() string {
	// fmt.Printf("ライン: %v", cd)
	null := `NOT NULL`
	if cd.Null {
		null = `NULL`
	}

	dflt := ``
	if 0 < len(cd.Default) {
		if cd.Default == "AUTO_INCREMENT" {
			dflt = fmt.Sprintf("%s", cd.Default)
		} else {
			dflt = fmt.Sprintf("DEFAULT %s", cd.Default)
		}
	}

	// key := ``
	// if 0 < len(cd.Key) {
	// 	key = cd.Key
	// }

	comment := ``
	if 0 < len(cd.Comment) {
		comment = fmt.Sprintf("comment '%s'", cd.Comment)
	}

	s := fmt.Sprintf("`%s` %s %s %s %s", cd.Name, cd.Type, null, dflt, comment)
	s = strings.Trim(s, ` `)
	rgx := regexp.MustCompile(" +")
	return rgx.ReplaceAllString(s, ` `)
}

type IndexDefinition struct {
	Columns  []string
	IsUnique bool
}

func (idx IndexDefinition) ToSQLStmt() string {
	if idx.IsUnique {
		return fmt.Sprintf("UNIQUE(`%s`)", strings.Join(idx.Columns, "`,`"))
	}

	return fmt.Sprintf("INDEX(`%s`)", strings.Join(idx.Columns, "`,`"))
}

type PrimaryKeysDefinition struct {
	Columns  []string
}

func (pk PrimaryKeysDefinition) ToSQLStmt() string {
	// fmt.Println("きてる？")
	return fmt.Sprintf("PRIMARY KEY (`%s`)", strings.Join(pk.Columns, "`,`"))
}

type ForeignKeysDefinition struct {
	Columns	[]string
	PreferenceTbl	string
	IndexNames	[]string
	OnUpdate	string
	OnDelete	string
}

func (fk ForeignKeysDefinition) ToSQLStmt() string {
	// fmt.Println(" FKきてる？")
	sql := fmt.Sprintf("FOREIGN KEY (`%s`)\n", strings.Join(fk.Columns, "`,`"))
	sql = sql + fmt.Sprintf("    REFERENCES `%s`(`%s`)\n", fk.PreferenceTbl, strings.Join(fk.IndexNames, "`,`"))
	sql = sql + fmt.Sprintf("    ON UPDATE %s ON DELETE %s", fk.OnUpdate, fk.OnDelete)
	return sql
}
