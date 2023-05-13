package util

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Query struct binder for default query param
type Query struct {
	Page    int    `json:"page"`
	Count   int    `json:"count"`
	Sort    string `json:"sort"`
	Filter  map[string]interface{}
	SetData map[string]interface{}
}

const (
	defaultPage  = 1
	defaultCount = 15
)

// NewQueryFilter initiate query with filter only
func NewQueryFilter(filter map[string]interface{}) *Query {
	return NewQuery(0, 0, "", filter, nil)
}

// NewQuery initiate query
func NewQuery(page, count int, sort string, filter map[string]interface{}, data map[string]interface{}) *Query {
	p := page
	if page < 1 {
		p = defaultPage
	}
	c := count
	if c < 1 {
		c = defaultCount
	}
	return &Query{
		Page:    p,
		Count:   c,
		Sort:    sort,
		Filter:  filter,
		SetData: data,
	}
}

// SetPage set value of page
func (q *Query) SetPage(i int) *Query {
	q.Page = i
	return q
}

// SetCount set value of count
func (q *Query) SetCount(i int) *Query {
	q.Count = i
	return q
}

// SetSort set value of sort
func (q *Query) SetSort(s string) *Query {
	q.Sort = s
	return q
}

// Operator string transalation
var Operator = map[string]string{
	"gt":      ">",      // elastic
	"lt":      "<",      // elastic
	"eq":      "=",      // elastic
	"ne":      "!=",     // elastic
	"gte":     ">=",     // elastic
	"lte":     "<=",     // elastic
	"like":    "like",   // elastic - text search
	"in":      "in",     // elastic
	"any":     "any",    // elastic / postgres
	"notin":   "not in", // elastic
	"null":    "is null",
	"notnull": "is not null",
}

// Where generate sql WHERE statement ,with format
//		key :"{columnName}{$operator}"
//		value : interface
// with default operator value "$eq"
// for example :
//     "amount$gte": 19200.00
// 	   "status": 1
// will be translated into sql format :
// 		WHERE amount >= 19200.00
//		AND status = 1
func (q *Query) Where() (string, []interface{}) {
	query := new(bytes.Buffer)
	var args []interface{}
	i := 0
	for k, v := range q.Filter {
		fields := strings.Split(k, "$")
		columnName := fields[0]
		if len(fields) < 2 {
			panic(fmt.Errorf("missing operator on query field %s", columnName))
		}
		opr := translateOperator(fields[1])
		isRequire := func(s string) bool {
			return s[len(s)-1:] == "!"
		}(fields[1])
		if i == 0 {
			isNull, _ := IsArgNil(v)
			if isRequire || !isNull {
				switch opr {
				case Operator["null"], Operator["notnull"]:
					query.WriteString(` WHERE ` + columnName + ` ` + opr + ` `)
				case Operator["like"]:
					tmpArgs, ok := v.(string)
					mm := `%` + strings.ToLower(tmpArgs) + `%`
					if ok {
						query.WriteString(` WHERE lower(` + columnName + `) ` + opr + ` ? `)
						args = append(args, mm)
					}
				case Operator["in"], Operator["notin"]:
					s := reflect.ValueOf(v)
					if s.Len() > 0 {
						if s.Kind() == reflect.Slice {
							var smt string
							for j := 0; j < s.Len(); j++ {
								smt += `?,`
								switch s.Index(j).Kind() {
								case reflect.Int:
									args = append(args, s.Index(j).Int())
								case reflect.Int32:
									args = append(args, s.Index(j).Int())
								case reflect.Int64:
									args = append(args, s.Index(j).Int())
								case reflect.Float32:
									args = append(args, s.Index(j).Float())
								case reflect.Float64:
									args = append(args, s.Index(j).Float())
								default:
									args = append(args, s.Index(j).String())
								}
							}
							query.WriteString(` WHERE ` + columnName + ` ` + opr + ` (` + smt[:len(smt)-1] + `)`)
						} else if s.Kind() == reflect.String {
							args = append(args, q.Escape(s.String()))
							query.WriteString(` WHERE ` + columnName + ` ` + opr + ` ( ? )`)
						}
					} else {
						args = append(args, nil)
						query.WriteString(` WHERE ` + columnName + ` ` + opr + ` ( ? ) `)
					}
					// query.WriteString(` WHERE ` + columnName + ` ` + opr + ` (?) `)
					// args = append(args, v)
				case Operator["any"]:
					s := reflect.ValueOf(v)
					if s.Kind() == reflect.Slice {
						if s.Len() > 0 {
							for j := 0; j < s.Len(); j++ {
								if j == 0 {
									query.WriteString(` WHERE '` + s.Index(j).String() + `' = ` + opr + `(` + columnName + `)`)
								} else {
									query.WriteString(` AND '` + s.Index(j).String() + `' = ` + opr + `(` + columnName + `)`)
								}
							}
						}
					}
				default:
					valueColumn, ok := isColumn(v)
					if ok {
						query.WriteString(` WHERE ` + columnName + ` ` + opr + ` ` + q.Escape(valueColumn) + ` `)
					} else {
						query.WriteString(` WHERE ` + columnName + ` ` + opr + ` ? `)
						args = append(args, v)
					}
				}
			} else {
				query.WriteString(` WHERE 1 = 1 `)
			}

		} else {
			isNull, _ := IsArgNil(v)
			if isRequire || !isNull {
				switch opr {
				case Operator["null"], Operator["notnull"]:
					query.WriteString(` AND ` + columnName + ` ` + opr + ` `)
				case Operator["like"]:
					tmpArgs, ok := v.(string)
					mm := `%` + strings.ToLower(tmpArgs) + `%`
					if ok {
						query.WriteString(` AND lower(` + columnName + `) ` + opr + ` ? `)
						args = append(args, mm)
					}
				case Operator["in"], Operator["notin"]:
					s := reflect.ValueOf(v)
					if s.Kind() == reflect.Slice {
						var smt string
						if s.Len() > 0 {
							for j := 0; j < s.Len(); j++ {
								smt += `?,`
								switch s.Index(j).Kind() {
								case reflect.Int:
									args = append(args, s.Index(j).Int())
								case reflect.Int32:
									args = append(args, s.Index(j).Int())
								case reflect.Int64:
									args = append(args, s.Index(j).Int())
								case reflect.Float32:
									args = append(args, s.Index(j).Float())
								case reflect.Float64:
									args = append(args, s.Index(j).Float())
								default:
									args = append(args, s.Index(j).String())
								}
							}
							query.WriteString(` AND ` + columnName + ` ` + opr + ` (` + smt[:len(smt)-1] + `)`)
						} else {
							args = append(args, nil)
							query.WriteString(` AND ` + columnName + ` ` + opr + ` ( ? )`)
						}
					} else if s.Kind() == reflect.String {
						args = append(args, q.Escape(s.String()))
						query.WriteString(` AND ` + columnName + ` ` + opr + ` ( ? )`)
					}
				case Operator["any"]:
					s := reflect.ValueOf(v)
					if s.Kind() == reflect.Slice {
						if s.Len() > 0 {
							for j := 0; j < s.Len(); j++ {
								query.WriteString(` AND '` + s.Index(j).String() + `' = ` + opr + `(` + columnName + `)`)
							}
						}
					}
				default:
					valueColumn, ok := isColumn(v)
					if ok {
						query.WriteString(` AND ` + columnName + ` ` + opr + ` ` + q.Escape(valueColumn) + ` `)
					} else {
						query.WriteString(` AND ` + columnName + ` ` + opr + ` ? `)
						args = append(args, v)
					}
				}
			} else {
				query.WriteString(` AND 1 = 1 `)
			}
		}
		i++
	}
	return query.String(), args
}

// Setter generate setter clause
func (q *Query) Setter() (string, []interface{}) {
	values := new(bytes.Buffer)
	var args []interface{}
	for k, v := range q.SetData {
		if v != nil {
			values.WriteString(k + ` = ? ,`)
			args = append(args, v)
		}
	}
	valuestr := values.String()

	return valuestr[:len(valuestr)-1], args
}

// Order generate string ordering query statement
func (q *Query) Order() string {
	if len(q.Sort) > 0 {
		field := strings.Split(q.Sort, ",")
		sort := `ORDER BY `
		for _, v := range field {
			sortType := func(str string) string {
				if strings.HasPrefix(str, "-") {
					return `DESC`
				}
				return `ASC`
			}
			sort += strings.TrimPrefix(v, "-") + ` ` + sortType(v) + `,`
		}
		return q.Escape(sort[:len(sort)-1])
	}
	return ``
}

// Limit generate limit and offset for pagination
func (q *Query) Limit() string {
	l := strconv.Itoa(int(q.Count))
	o := strconv.Itoa((int(q.Page - 1)) * int(q.Count))
	return ` LIMIT ` + l + ` OFFSET ` + o
}

func (q *Query) Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: // Must be escaped for 'mysql'
			escape = '0'
		case '\n': // Must be escaped for logs
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': // Better safe than sorry
			escape = '"'
		case '\032': //This gives problems on Win32
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

func translateOperator(s string) string {
	str := strings.Trim(s, "!")
	operator := Operator[strings.ToLower(str)]
	if operator == "" {
		return Operator["eq"]
	}
	return operator
}

func IsArgNil(i interface{}) (bool, reflect.Kind) {
	r := reflect.ValueOf(i)
	switch r.Kind() {
	case reflect.Slice:
		return r.Len() == 0, reflect.Slice
	case reflect.String:
		return r.String() == "", reflect.String
	case reflect.Int:
		return r.Int() == 0, reflect.Int
	case reflect.Int32:
		return r.Int() == 0, reflect.Int32
	case reflect.Int64:
		return r.Int() == 0, reflect.Int64
	case reflect.Float32:
		return r.Float() == 0, reflect.Float32
	case reflect.Float64:
		return r.Float() == 0, reflect.Float64
	default:
		return false, reflect.String
	}
}

func isColumn(i interface{}) (string, bool) {
	col, ok := i.(string)
	if ok && strings.Contains(col, ":") {
		split := strings.Split(col, ":")
		if split[0] == "column" {
			return split[1], ok
		}
	}
	return col, false
}
