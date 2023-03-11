package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tag        string = "validate"
	ruleMin    string = "min"
	ruleMax    string = "max"
	ruleIn     string = "in"
	ruleLen    string = "len"
	ruleRegExp string = "regexp"
)

var (
	ErrWrongType         = errors.New("invalid type")
	ErrWrongValueLessMin = errors.New("wrong value less min")
	ErrWrongValueMoreMax = errors.New("wrong value more max")
	ErrWrongValueNotIn   = errors.New("wrong value not in")
	ErrWrongValueLength  = errors.New("wrong value length")
	ErrWrongValueRegExp  = errors.New("wrong value regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	err := strings.Builder{}

	for i, vv := range v {
		err.WriteString(fmt.Sprintf("%s: %s", vv.Field, vv.Err.Error()))

		if i < len(v)-1 {
			err.WriteString(", ")
		}
	}

	return err.String()
}

func Validate(v interface{}) error {
	var errors ValidationErrors

	value := reflect.ValueOf(v)

	if value.Kind().String() != "struct" {
		return ErrWrongType
	}

	for i := 0; i < value.NumField(); i++ {
		switch value.Field(i).Kind().String() {
		case "int":
			rules, ok := value.Type().Field(i).Tag.Lookup(tag)
			if !ok {
				continue
			}

			errorsValidation, err := validateInt(int(value.Field(i).Int()), value.Type().Field(i).Name, rules)
			if err != nil {
				return err
			}

			errors = append(errors, errorsValidation...)

		case "string":
			rules, ok := value.Type().Field(i).Tag.Lookup(tag)
			if !ok {
				continue
			}

			errorsValidation, err := validateString(value.Field(i).String(), value.Type().Field(i).Name, rules)
			if err != nil {
				return err
			}

			errors = append(errors, errorsValidation...)

		case "slice":
			if value.Field(i).CanInterface() {
				intSlice, ok := value.Field(i).Interface().([]int)

				if ok {
					rules, ok := value.Type().Field(i).Tag.Lookup(tag)
					if !ok {
						continue
					}

					errorsValidation, err := validateIntSlice(intSlice, value.Type().Field(i).Name, rules)
					if err != nil {
						return err
					}

					errors = append(errors, errorsValidation...)
				}

				stringSlice, ok := value.Field(i).Interface().([]string)

				if ok {
					rules, ok := value.Type().Field(i).Tag.Lookup(tag)
					if !ok {
						continue
					}

					errorsValidation, err := validateStringSlice(stringSlice, value.Type().Field(i).Name, rules)
					if err != nil {
						return err
					}

					errors = append(errors, errorsValidation...)
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func validateStringSlice(value []string, name, rules string) (errors ValidationErrors, err error) {
	for _, v := range value {
		errorsValidation, err := validateString(v, name, rules)
		if err != nil {
			return nil, err
		}

		if errorsValidation != nil {
			e := newValidationError(name, fmt.Errorf("%w value %v", errorsValidation[0].Err, v))

			errors = append(errors, e)
		}
	}

	return errors, nil
}

func validateIntSlice(value []int, name, rules string) (errors ValidationErrors, err error) {
	for _, v := range value {
		errorsValidation, err := validateInt(v, name, rules)
		if err != nil {
			return nil, err
		}

		if errorsValidation != nil {
			e := newValidationError(name, fmt.Errorf("%w value %v", errorsValidation[0].Err, v))

			errors = append(errors, e)
		}
	}

	return errors, nil
}

func validateString(value, name, rules string) (errors ValidationErrors, err error) {
	for _, rulesFields := range strings.Split(rules, "|") {
		rule := strings.Split(rulesFields, ":")

		switch rule[0] {
		case ruleLen:
			length, err := strconv.Atoi(rule[1])
			if err != nil {
				return nil, err
			}

			if len(value) != length {
				errors = append(errors, newValidationError(name, ErrWrongValueLength))
			}

		case ruleRegExp:
			ok, err := regexp.Match(rule[1], []byte(value))
			if err != nil {
				return nil, err
			}

			if !ok {
				errors = append(errors, newValidationError(name, ErrWrongValueRegExp))
			}

		case ruleIn:
			if !strings.Contains(rule[1], value) {
				errors = append(errors, newValidationError(name, ErrWrongValueNotIn))
			}
		}
	}

	return errors, nil
}

func validateInt(value int, name, rules string) (errors ValidationErrors, err error) {
	for _, rulesFields := range strings.Split(rules, "|") {
		rule := strings.Split(rulesFields, ":")

		switch rule[0] {
		case ruleMin:
			min, err := strconv.Atoi(rule[1])
			if err != nil {
				return nil, err
			}

			if value < min {
				errors = append(errors, newValidationError(name, ErrWrongValueLessMin))
			}

		case ruleMax:
			max, err := strconv.Atoi(rule[1])
			if err != nil {
				return nil, err
			}

			if value > max {
				errors = append(errors, newValidationError(name, ErrWrongValueMoreMax))
			}

		case ruleIn:
			ruleValue := strings.Split(rule[1], ",")

			ok, err := contains(ruleValue, value)
			if err != nil {
				return nil, err
			}

			if !ok {
				errors = append(errors, newValidationError(name, ErrWrongValueNotIn))
			}
		}
	}

	return errors, nil
}

func newValidationError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func contains(arr []string, value int) (bool, error) {
	for _, str := range arr {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return false, err
		}

		if value == intValue {
			return true, nil
		}
	}

	return false, nil
}
