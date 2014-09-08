package hat

import "strings"

type fieldFilter struct {
	Fields []string
}

func NewFieldFilter(csv string) *fieldFilter {
	if len(csv) == 0 {
		return &fieldFilter{[]string{"*"}}
	}
	raw := strings.Split(csv, ",")
	println("Raw len: ", len(raw))
	compiled := []string{}
	for _, r := range raw {
		multiplex := [][]string{}
		parts := strings.Split(r, ".")
		for _, p := range parts {
			if p[0:1] == "[" {
				p = strings.Trim(p, "[]")
				multiplex = append(multiplex, strings.Split(p, ":"))
			} else {
				multiplex = append(multiplex, []string{p})
			}
		}
		compiled = append(compiled, mp(multiplex)...)
	}
	return &fieldFilter{Fields: compiled}
}

func (ff *fieldFilter) Enter(fieldName string) *fieldFilter {
	newFields := []string{}
	for _, filter := range ff.Fields {
		parts := strings.SplitN(filter, ".", 2)
		if len(parts) == 2 && parts[0] == fieldName {
			newFields = append(newFields, parts[1])
		}
	}
	return &fieldFilter{newFields}
}

func (ff *fieldFilter) Allows(fieldName string) bool {
	for _, filter := range ff.Fields {
		parts := strings.SplitN(filter, ".", 2)
		if parts[0] == fieldName || parts[0] == "*" {
			return true
		}
	}
	return false
}

func (ff *fieldFilter) Filter(v interface{}) (interface{}, error) {
	if len(ff.Fields) == 0 {
		// No more filters, so return whole value
		return v, nil
	}
	slice, isSlice := v.([]interface{})
	if isSlice {
		return ff.FilterSlice(slice)
	}
	return ff.FilterNonSlice(v)
}

func (ff *fieldFilter) FilterSlice(slice []interface{}) (interface{}, error) {
	val := make([]interface{}, len(slice))
	for i, v := range slice {
		v, err := ff.Filter(v)
		if err != nil {
			return nil, err
		}
		val[i] = v
	}
	return val, nil
}

func (ff *fieldFilter) FilterNonSlice(v interface{}) (interface{}, error) {
	m, err := toSmap(v)
	if err != nil {
		return nil, err
	}
	fm := smap{}
	for k, v := range m {
		if ff.Allows(k) {
			result, err := ff.Enter(k).Filter(v)
			if err != nil {
				return nil, err
			}
			fm[k] = result
		}
	}
	return fm, nil
}

// mp multiplies a 2d string array into a single array using dots as separators.
// e.g. mp([[a,b,c],[d,e]]) == [a.d, b.d, c.d, a.e, b.e, c.e]
func mp(a [][]string) []string {
	out := a[0]
	for i := 1; i < len(a); i++ {
		out = multiply(out, a[i])
	}
	return out
}

func multiply(a []string, b []string) []string {
	if len(a) == 0 {
		a = []string{""}
	}
	if len(b) == 0 {
		b = []string{""}
	}
	out := []string{}
	for _, s := range a {
		for _, s1 := range b {
			out = append(out, s+"."+s1)
		}
	}
	return out
}
