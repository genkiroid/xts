package xts

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// A StructureField represents fields of column information.
type StructureField struct {
	Field   string `xml:"Field,attr"`
	Type    string `xml:"Type,attr"`
	Null    string `xml:"Null,attr"`
	Key     string `xml:"Key,attr"`
	Default string `xml:"Default,attr"`
	Extra   string `xml:"Extra,attr"`
	Comment string `xml:"Comment,attr"`
}

// A TableStructure represents a table structure information.
type TableStructure struct {
	Name   string           `xml:"name,attr"`
	Fields []StructureField `xml:"field"`
}

// A DataField represents a column data.
type DataField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
	Null  string `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
}

// A DataRow represents a mysql data record.
type DataRow struct {
	Fields []DataField `xml:"field"`
}

// A TableData represents a record set of a table.
type TableData struct {
	Name     string    `xml:"name,attr"`
	DataRows []DataRow `xml:"row"`
}

// A Database represents a database.
type Database struct {
	Name            string           `xml:"name,attr"`
	TableStructures []TableStructure `xml:"table_structure"`
	TableData       []TableData      `xml:"table_data"`
}

// A MySQLXMLDump represents an unmarshaled of a MySQL xml dump file.
type MySQLXMLDump struct {
	XMLName  xml.Name `xml:"mysqldump"`
	Database Database `xml:"database"`
}

// NewMySQLXMLDump returns a new NewMySQLXMLDump.
func NewMySQLXMLDump(raw []byte) (MySQLXMLDump, error) {
	v := MySQLXMLDump{}
	err := xml.Unmarshal(raw, &v)
	return v, err
}

// A Columns represents columns of a table to build sql.
type Columns []string

// String returns comma separated string of column names.
func (c Columns) String() string {
	return strings.Join(c, ", ")
}

// A Value represents a set of information of a column value to build sql.
type Value struct {
	Name  string
	Type  string
	Value string
	Null  string
}

// Yaml returns string as yaml format for a value.
func (v Value) Yaml() string {
	if v.Null == "true" {
		return fmt.Sprintf("%s: NULL\n", v.Name)
	}
	if v.Value == "" {
		return fmt.Sprintf("%s: \"\"\n", v.Name)
	}
	if strings.Contains(v.Value, "\n") {
		return fmt.Sprintf("%s: \"%s\"\n", v.Name, strings.ReplaceAll(v.Value, "\n", "\\n"))
	}
	return fmt.Sprintf("%s: %s\n", v.Name, v.Value)
}

// A Values represents slice of Value.
type Values []Value

// String returns comma separated string of values for a record.
func (v Values) String() string {
	var s []string
	for _, value := range v {
		if value.Null == "true" {
			s = append(s, "NULL")
			continue
		}
		t := strings.Split(value.Type, "(")[0]
		switch t {
		case "bit", "tinyint", "bool", "boolean", "smallint", "mediumint", "int", "integer", "bigint", "decimal", "dec", "float", "double":
			s = append(s, value.Value)
		case "date", "datetime", "timestamp", "time", "year", "char", "varchar", "binary", "varbinary", "tinyblob", "tinytext", "blob", "text", "midiumtext", "longblob", "longtext", "enum", "set":
			s = append(s, fmt.Sprintf("'%s'", strings.ReplaceAll(value.Value, "\n", "\\n")))
		default:
		}
	}
	return strings.Join(s, ", ")
}

// Yaml returns string as yaml format for values.
func (v Values) Yaml() string {
	var s []string
	for _, value := range v {
		s = append(s, value.Yaml())
	}
	return fmt.Sprintf("- %s", strings.Join(s, "  "))
}

// A Rows represents slice of Values.
type Rows []Values

// String returns comma separated string of values for record set.
func (r Rows) String() string {
	var s []string
	for _, values := range r {
		s = append(s, fmt.Sprintf("(%s)", values))
	}
	return strings.Join(s, ", ")
}

// Yaml returns string as yaml format for data rows.
func (r Rows) Yaml() string {
	var s []string
	for _, values := range r {
		s = append(s, values.Yaml())
	}
	return strings.Join(s, "\n")
}

// A Table represents a table to build sql.
type Table struct {
	Name    string
	Columns Columns
	Rows    Rows
}

// A Sql has slice of Tables to build sql.
type Sql struct {
	Tables []Table
	out    string
}

// NewSql returns a new Sql.
func NewSql(x MySQLXMLDump, out string) Sql {
	var sql Sql
	sql.out = out

	tableStructures := x.Database.TableStructures
	for i, table := range tableStructures {
		var columns []string
		for _, f := range table.Fields {
			columns = append(columns, f.Field)
		}

		var rows Rows
		for _, row := range x.Database.TableData[i].DataRows {
			values := Values{}
			for j, field := range row.Fields {
				values = append(values, Value{field.Name, tableStructures[i].Fields[j].Type, field.Value, field.Null})
			}
			rows = append(rows, values)
		}
		sql.Tables = append(sql.Tables, Table{table.Name, columns, rows})
	}
	return sql
}

// InsertStmt output insert sql statement built from xml dump file to specified out.
func (sql Sql) InsertStmt() error {
	for _, table := range sql.Tables {
		writer, err := sql.writer(table.Name + ".sql")
		if err != nil {
			return fmt.Errorf("sql.writer(%s) returned error: %s", table.Name, err)
		}
		fmt.Fprintf(writer, "INSERT INTO %s (%s) VALUES %s;\n", table.Name, table.Columns, table.Rows)
	}
	return nil
}

// Yaml output yaml format built from xml dump file to specified out.
func (sql Sql) Yaml() error {
	for _, table := range sql.Tables {
		writer, err := sql.writer(table.Name + ".yml")
		if err != nil {
			return fmt.Errorf("sql.writer(%s) returned error: %s", table.Name, err)
		}
		fmt.Fprintf(writer, "# %s\n%s\n", table.Name, table.Rows.Yaml())
	}
	return nil
}

// writer returns io.Writer by specified out.
func (sql Sql) writer(f string) (io.Writer, error) {
	path := sql.out
	if path == "" {
		return os.Stdout, nil
	}

	var err error
	if path == "--" {
		path, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("os.Getwd() returned error: %s\n", err)
		}
	}

	path = filepath.Join(path, f)

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("os.Create(%s) returned error: %s\n", path, err)
	}
	return file, nil
}
