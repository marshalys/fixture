package loaders

import (
	"fmt"
	"strings"
)

type LoadContent struct {
	Version string                   `yaml:"version" json:"version"`
	Table   string                   `yaml:"table" json:"table"`
	Rows    []map[string]interface{} `yaml:"rows" json:"rows"`
}

var sqlTemplate = `INSERT INTO %s (%s)
VALUES 
 %s;
`

func toInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func getFieldText(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "NULL"
	case bool:
		return fmt.Sprintf("b'%v'", toInt(v))
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func genSQL(content *LoadContent) string {
	if content == nil || content.Table == "" {
		return ""
	}

	if len(content.Rows) == 0 {
		return ""
	}

	columns := make([]string, 0)
	for k := range content.Rows[0] {
		columns = append(columns, k)
	}

	vals := make([]string, 0)
	for _, row := range content.Rows {
		fields := make([]string, 0)
		for _, col := range columns {
			if val, ok := row[col]; ok {
				fields = append(fields, getFieldText(val))
			} else {
				panic(fmt.Sprintf("fixture.loaders: incosistent column found '%s'", col))
			}
		}

		vals = append(vals, "("+strings.Join(fields, ", ")+")")
	}

	for i, col := range columns {
		columns[i] = quote(col)
	}

	exp := fmt.Sprintf(sqlTemplate, quote(content.Table), strings.Join(columns, ", "), strings.Join(vals, ",\n"))
	return exp
}

func quote(col string) string {
	if col != "" {
		col = fmt.Sprintf("`%s`", col)
	}
	return col
}
