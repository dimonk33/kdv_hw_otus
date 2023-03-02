package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ErrNotStruct          = errors.New("variable not a struct")
	ErrNoTags             = errors.New("no Tags. Skip field")
	ErrNoValidTag         = errors.New("no Tag validate. Skip field")
	ErrNotIntSlice        = errors.New("error get []int")
	ErrNotStringSlice     = errors.New("error get []string")
	ErrWrongMinCond       = errors.New("error int min condition")
	ErrWrongMaxCond       = errors.New("error int max condition")
	ErrWrongArrayCond     = errors.New("error array condition")
	ErrWrongLenCond       = errors.New("error string len condition")
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

func (v ValidationErrors) Error() string {
	var errMerge string
	for i, e := range v {
		errMerge += strconv.Itoa(i+1) + ") " + e.Field + ":" + e.Err.Error() + "\n"
	}
	return errMerge
}

func Validate(v interface{}) error {
	outErrors := ValidationErrors{}
	types := reflect.TypeOf(v)
	values := reflect.ValueOf(v)
	if types.Kind().String() != "struct" {
		return ErrNotStruct
	}
	for i := 0; i < types.NumField(); i++ {
		varName := types.Field(i).Name
		varType := types.Field(i).Type
		varTag := types.Field(i).Tag
		varValue := values.Field(i)
		fmt.Printf("%v %v %v %v\n", varName, varType, varTag, varValue)

		if len(varTag) == 0 {
			err := ValidationError{
				Field: varName,
				Err:   ErrNoTags,
			}
			outErrors = append(outErrors, err)
			continue
		}
		validateCondition, ok := varTag.Lookup("validate")
		if !ok {
			err := ValidationError{
				Field: varName,
				Err:   ErrNoValidTag,
			}
			outErrors = append(outErrors, err)
			continue
		}

		switch varType.Kind().String() {
		case "int":
			err := validateInt(varValue.Int(), validateCondition)
			if err != nil {
				outErrors = append(outErrors, ValidationError{
					Field: varName,
					Err:   fmt.Errorf(varType.Kind().String()+"validate: %w", err),
				})
			}
		case "[]int":
			slice, ok := varValue.Interface().([]int)
			if !ok {
				err := ValidationError{
					Field: varName,
					Err:   ErrNotIntSlice,
				}
				outErrors = append(outErrors, err)
			}
			err := validateIntSlice(slice, validateCondition)
			if err != nil {
				outErrors = append(outErrors, ValidationError{
					Field: varName,
					Err:   fmt.Errorf(varType.Kind().String()+"validate: %w", err),
				})
			}
		case "string":
			err := validateString(varValue.String(), validateCondition)
			if err != nil {
				outErrors = append(outErrors, ValidationError{
					Field: varName,
					Err:   fmt.Errorf(varType.Kind().String()+"validate: %w", err),
				})
			}
		case "[]string":
			slice, ok := varValue.Interface().([]string)
			if !ok {
				outErrors = append(outErrors, ValidationError{
					Field: varName,
					Err:   ErrNotStringSlice,
				})
			}
			err := validateStringSlice(slice, validateCondition)
			if err != nil {
				outErrors = append(outErrors, ValidationError{
					Field: varName,
					Err:   fmt.Errorf(varType.Kind().String()+"validate: %w", err),
				})
			}
		}
	}
	return outErrors
}

func validateInt(value int64, conditionStr string) error {
	fmt.Printf("validate int: %d %v\n", value, conditionStr)
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
			var numArr []int
			for _, numStr := range numStrArr {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return ErrWrongArrayCond
				}
				numArr = append(numArr, num)
			}
			if !slices.Contains(numArr, int(value)) {
				return ErrValidateArrayCond
			}
		}
	}

	return nil
}

func validateIntSlice(value []int, conditionStr string) error {
	fmt.Printf("validate []int: %v %v\n", value, conditionStr)
	for _, v := range value {
		if err := validateInt(int64(v), conditionStr); err != nil {
			return err
		}
	}
	return nil
}

func validateString(value string, conditionStr string) error {
	fmt.Printf("validate string: %s %v\n", value, conditionStr)
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

func validateStringSlice(value []string, conditionStr string) error {
	fmt.Printf("validate []string: %v %v\n", value, conditionStr)
	for _, v := range value {
		if err := validateString(v, conditionStr); err != nil {
			return err
		}
	}
	return nil
}
