package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ErrNotStruct          = errors.New("variable not a struct")
	ErrNotIntSlice        = errors.New("error get []int")
	ErrNotStringSlice     = errors.New("error get []string")
	ErrWrongMinCond       = errors.New("error int min condition")
	ErrWrongMaxCond       = errors.New("error int max condition")
	ErrWrongArrayCond     = errors.New("error array condition")
	ErrWrongLenCond       = errors.New("error string len condition")
	ErrWrongType          = errors.New("error not support field")
	ErrValidateMinCond    = errors.New("value < min")
	ErrValidateMaxCond    = errors.New("value > max")
	ErrValidateArrayCond  = errors.New("value not in array condition")
	ErrValidateLenCond    = errors.New(" string len != len")
	ErrValidateRegexpCond = errors.New("value not match regexp condition")
)

const (
	tagMin    = "min:"
	tagMax    = "max:"
	tagIn     = "in:"
	tagLen    = "len:"
	tagRegexp = "regexp:"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return v.Err.Error()
}

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for i, e := range v {
		builder.WriteString(strconv.Itoa(i+1) + ") " + e.Field + ":" + e.Err.Error() + "\n")
	}
	return builder.String()
}

func Validate(v interface{}) error {
	outErrors := ValidationErrors{}
	var exit bool
	var err error
	types := reflect.TypeOf(v)
	values := reflect.ValueOf(v)
	if types.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < types.NumField(); i++ {
		varName := types.Field(i).Name
		varType := types.Field(i).Type
		varTag := types.Field(i).Tag
		varValue := values.Field(i)

		if !varValue.CanSet() || len(varTag) == 0 {
			continue
		}

		validateCondition, ok := varTag.Lookup("validate")
		if !ok {
			continue
		}

		switch varType.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int8, reflect.Int64, reflect.Int16,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = validateInt(varValue.Int(), validateCondition)
			outErrors, exit = addValidateError(varName, err, outErrors)
		case reflect.String:
			err = validateString(varValue.String(), validateCondition)
			outErrors, exit = addValidateError(varName, err, outErrors)
		case reflect.Slice:
			err = validateSlice(varType, varValue, validateCondition)
			outErrors, exit = addValidateError(varName, err, outErrors)
		default:
			return ErrWrongType
		}

		if exit {
			return err
		}
	}
	return outErrors
}

func validateInt(value int64, conditionStr string) error {
	conditions := strings.Split(conditionStr, "|")
	for _, condition := range conditions {
		switch {
		case strings.Contains(condition, tagMin):
			numStr := strings.ReplaceAll(condition, tagMin, "")
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return ErrWrongMinCond
			}
			if value < int64(num) {
				return ErrValidateMinCond
			}
		case strings.Contains(condition, tagMax):
			numStr := strings.ReplaceAll(condition, tagMax, "")
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return ErrWrongMaxCond
			}
			if value > int64(num) {
				return ErrValidateMaxCond
			}
		case strings.Contains(condition, tagIn):
			numArrStr := strings.ReplaceAll(condition, tagIn, "")
			numStrArr := strings.Split(numArrStr, ",")
			isFind := false
			for _, numStr := range numStrArr {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return ErrWrongArrayCond
				}
				if value == int64(num) {
					isFind = true
					break
				}
			}
			if !isFind {
				return ErrValidateArrayCond
			}
		}
	}

	return nil
}

func validateString(value string, conditionStr string) error {
	conditions := strings.Split(conditionStr, "|")

	for _, condition := range conditions {
		switch {
		case strings.Contains(condition, tagLen):
			numStr := strings.ReplaceAll(condition, tagLen, "")
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return ErrWrongLenCond
			}
			if len(value) != num {
				return ErrValidateLenCond
			}
		case strings.Contains(condition, tagRegexp):
			re := strings.ReplaceAll(condition, tagRegexp, "")
			match, _ := regexp.MatchString(re, value)
			if !match {
				return ErrValidateRegexpCond
			}
		case strings.Contains(condition, tagIn):
			arrStr := strings.ReplaceAll(condition, tagIn, "")
			strArr := strings.Split(arrStr, ",")
			if !slices.Contains(strArr, value) {
				return ErrValidateArrayCond
			}
		}
	}

	return nil
}

func validateSlice(
	varType reflect.Type,
	varValue reflect.Value,
	validateCondition string,
) error {
	switch varType.Elem().Kind() {
	case reflect.Int:
		slice, ok := varValue.Interface().([]int)
		if !ok {
			return ErrNotIntSlice
		}
		err := validateIntSlice(slice, validateCondition)
		if err != nil {
			return err
		}
	case reflect.String:
		slice, ok := varValue.Interface().([]string)
		if !ok {
			return ErrNotStringSlice
		}
		err := validateStringSlice(slice, validateCondition)
		if err != nil {
			return err
		}
	default:
		return ErrWrongType
	}
	return nil
}

func validateIntSlice(value []int, conditionStr string) error {
	for _, v := range value {
		if err := validateInt(int64(v), conditionStr); err != nil {
			return err
		}
	}
	return nil
}

func validateStringSlice(value []string, conditionStr string) error {
	for _, v := range value {
		if err := validateString(v, conditionStr); err != nil {
			return err
		}
	}
	return nil
}

func addValidateError(varName string, err error, outErrors ValidationErrors) (ValidationErrors, bool) {
	if !(errors.Is(err, ErrWrongMinCond) ||
		errors.Is(err, ErrWrongMaxCond) ||
		errors.Is(err, ErrWrongArrayCond) ||
		errors.Is(err, ErrWrongLenCond)) {
		outErrors = append(outErrors, ValidationError{
			Field: varName,
			Err:   err,
		})
		return outErrors, false
	}
	return outErrors, true
}
