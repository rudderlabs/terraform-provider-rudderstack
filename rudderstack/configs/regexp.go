package configs

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func StringMatchesRegexp(rs string) schema.SchemaValidateDiagFunc {
	r, err := regexp.Compile(rs)
	if err != nil {
		panic(err)
	}

	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("value for %q is not of type string", k)}
		}
		if ok := r.MatchString(v); !ok {
			return nil, []error{fmt.Errorf("value for %q does not match regular expression %q", k, r)}
		}

		return nil, nil
	})
}
