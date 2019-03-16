package main

import (
	"fmt"
	"../../sam"
	"../../sqldef"
	"os"
	"strings"
)

func onTableColumns(tbl *sqldef.TableDefinition, columns map[string]string) error {

	nullable := false
	if columns[`NULL`] == `YES` {
		nullable = true
	}
	// fmt.Printf(" >> %v\n", columns)
	// fmt.Printf(" >> %v\n\n", columns[`論理名`])


	colDef := sqldef.ColumnDefinition{
		// -               Name:    columns[`name`],
		// -               Type:    columns[`type`],
		// +               Desc:    columns[`論理名`],
		// +               Name:    columns[`物理名`],
		// +               Type:    columns[`型`],
		// 				Null:    nullable,
		// -               Default: columns[`default`],
		// -               Key:     columns[`key`],
		// -               Comment: columns[`comment`],
		// +               Default: columns[`DEFAULT`],
		// +               Key:     columns[`KEY`],
		// +               Comment: columns[`COMMENT`],		Desc:	 columns[`論理名`],
		Desc:    columns[`論理名`],
		Name:    columns[`物理名`],
		Type:    columns[`型`],
		Null:    nullable,
		Default: columns[`DEFAULT`],
		Key:     columns[`KEY`],
		Comment: columns[`COMMENT`],
	}
	// fmt.Printf(" >> %v こいつ？", colDef)
	tbl.Columns = append(tbl.Columns, colDef)
	return nil
}

func onTableIndexes(tbl *sqldef.TableDefinition, columns map[string]string) error {
	idxDef := sqldef.IndexDefinition{
		Columns:  strings.Split(columns[`columns`], `,`),
		IsUnique: columns[`unique`] == `YES`,
	}

	tbl.Indexes = append(tbl.Indexes, idxDef)
	return nil
}

func onTablePraimaryKeys(tbl *sqldef.TableDefinition, columns map[string]string) error {
	// fmt.Println(" -- ここきてる？ ")
	pkDef := sqldef.PrimaryKeysDefinition{
		Columns:  strings.Split(columns[`columns`], `,`),
	}

	tbl.PrimaryKeys = append(tbl.PrimaryKeys, pkDef)
	return nil
}

func onTableForeignKeys(tbl *sqldef.TableDefinition, columns map[string]string) error {
	// fmt.Println(" -- ここきてる？!!! ")
	fkDef := sqldef.ForeignKeysDefinition{
		Columns:  strings.Split(columns[`columns`], `,`),
		PreferenceTbl:  columns[`preference_tbl`],
		IndexNames:  strings.Split(columns[`index_names`], `,`),
		OnUpdate:  columns[`constraint_update`],
		OnDelete:  columns[`constraint_delete`],
	}

	tbl.ForeignKeys = append(tbl.ForeignKeys, fkDef)
	return nil
}

func main() {
	// fmt.Println("デバッグ")
	mdPath := os.Args[1]

	sqlDef := sqldef.SQLDefinition{}
	currentTableIdx := -1

	isColumnMode := false
	isIndexesMode := false
	isPKMode := false
	isFKMode := false

	sm := sam.SamParser{
		OnOneLines: map[string]func(line string) error{
			"#": func(line string) error {
				sqlDef.DatabaseName = line
				return nil
			},
			"##": func(line string) error {
				tblDef := sqldef.TableDefinition{Name: line}
				sqlDef.Tables = append(sqlDef.Tables, tblDef)
				currentTableIdx = currentTableIdx + 1
				return nil
			},
			"###": func(line string) error {
				isColumnMode = false
				isIndexesMode = false
				isPKMode = false
				isFKMode = false

				// fmt.Printf(" >> ### のとき %v\n", line)

				if line == "columns" {
					isColumnMode = true
					return nil
				}

				// インデックス文
				if line == "Indexes" {
					isIndexesMode = true
					return nil
				}

				// PK文
				if line == "PrimaryKeys" {
					isPKMode = true
					return nil
				}

				// FK文
				if line == "ForeignKeys" {
					isFKMode = true
					return nil
				}

				return fmt.Errorf("Unkown ### [%s]", line)
			},
		},
		OnTable: func(columns map[string]string) error {
			tblDef := &sqlDef.Tables[currentTableIdx]

			// fmt.Printf(" >>>>> columns, %v\n", columns)

			// キー名、値のトリム
			trimedColumns := make(map[string]string)
			for key, val := range columns {
				trimedKey := strings.Trim(key, " ")
				trimedVal := strings.Trim(val, " ")
				trimedColumns[trimedKey] = trimedVal
			}
			// fmt.Printf(" >>>>> trimedColumns, %v\n", len(trimedColumns))
			// fmt.Printf(" >>>>> isPKMode, %v\n", isPKMode)

			if isColumnMode {
				return onTableColumns(tblDef, trimedColumns)
			}

			if isIndexesMode {
				return onTableIndexes(tblDef, trimedColumns)
			}

			if isPKMode {
				return onTablePraimaryKeys(tblDef, trimedColumns)
			}

			if isFKMode {
				return onTableForeignKeys(tblDef, trimedColumns)
			}

			return nil
		},
	}

	// fmt.Printf(" >>>> %v\n", sm)

	if err := sm.Start(mdPath); err != nil {
		panic(err)
	}

	fmt.Println(sqlDef.ToSQLStmt())
}
