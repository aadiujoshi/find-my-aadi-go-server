package findmyaadigoserver

import (
    "fmt"
    "strconv"
)

// Int64ToStr converts int64 to string
func Int64ToStr(i int64) string {
    return strconv.FormatInt(i, 10)
}

// StrToInt64 converts string to int64
func StrToInt64(s string) (int64, error) {
    return strconv.ParseInt(s, 10, 64)
}

// Float64ToStr converts float64 to string (Java-style, SQL-friendly)
func Float64ToStr(f float64) string {
    // 'g' = compact, -1 = full precision, 64 = float64
    return strconv.FormatFloat(f, 'g', -1, 64)
}

// StrToFloat64 converts string to float64
func StrToFloat64(s string) (float64, error) {
    return strconv.ParseFloat(s, 64)
}

// Optional helper for SQL-friendly quoting
func SQLValue(v interface{}) string {
    switch val := v.(type) {
    case int64:
        return Int64ToStr(val)
    case float64:
        return Float64ToStr(val)
    case string:
        return fmt.Sprintf("'%s'", val)
    default:
        return fmt.Sprintf("'%v'", val)
    }
}
