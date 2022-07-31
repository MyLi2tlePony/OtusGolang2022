package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
)

type ValidationError struct {
	Field string
	Err   error
}

type Validator struct {
	key   string
	value string
}

type ValidationErrors []ValidationError

var (
	regexpGetValidators       = regexp.MustCompile(`(regexp:\^.+\$)|([^|]+)`)
	regexpGetValidatorKey     = regexp.MustCompile(`^[^:]+`)
	regexpGetValidatorValue   = regexp.MustCompile(`[^:]+$`)
	regexpParseValidatorValue = regexp.MustCompile(`[^,]+`)

	ErrInterfaceNotPointer = errors.New("interface must be a pointer")
	ErrInterfaceNotStruct  = errors.New("interface must be a struct")

	ErrNotImplementedType = errors.New("not implemented type")
	ErrNotImplementedTag  = errors.New("not implemented tag")

	ErrStringRegexp = errors.New("error string regexp")
	ErrStringLen    = errors.New("error string len")
	ErrStringIn     = errors.New("error string in")

	ErrIntMin = errors.New("error int min")
	ErrIntMax = errors.New("error int max")
	ErrIntIn  = errors.New("error int in")
)

const tagName = "validate"

func (v ValidationErrors) Error() string {
	resultError := ""

	for _, validationError := range v {
		resultError += validationError.Field + ": " + validationError.Err.Error() + "\n"
	}

	return resultError
}

func Validate(item interface{}) error {
	if reflect.ValueOf(item).Kind() != reflect.Struct {
		return ErrInterfaceNotStruct
	}

	value := reflect.ValueOf(item).Elem()
	if !value.CanAddr() {
		return ErrInterfaceNotPointer
	}

	validationErrors := make(ValidationErrors, 0)

	for fieldID := 0; fieldID < value.NumField(); fieldID++ {
		structField := value.Type().Field(fieldID)

		if tag, ok := structField.Tag.Lookup(tagName); ok {
			validators := getValidators(tag)
			kind := value.Field(fieldID).Kind()
			field := value.Field(fieldID)

			if err := validateField(validators, kind, field); err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: structField.Name,
					Err:   err,
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func getValidators(tag string) []Validator {
	validators := make([]Validator, 0)
	v := regexpGetValidators.FindAllString(tag, -1)

	for _, value := range v {
		validators = append(validators, Validator{
			key:   regexpGetValidatorKey.FindString(value),
			value: regexpGetValidatorValue.FindString(value),
		})
	}

	return validators
}

func validateField(validators []Validator, kind reflect.Kind, field reflect.Value) error {
	switch kind { // nolint:exhaustive
	case reflect.String:
		value := field.String()

		if err := validateString(value, validators); err != nil {
			return err
		}

	case reflect.Int:
		value := field.Int()

		if err := validateInt(value, validators); err != nil {
			return err
		}

	case reflect.Slice:
		if err := validateSlice(field, validators); err != nil {
			return err
		}

	default:
		return ErrNotImplementedType
	}

	return nil
}

func validateSlice(field reflect.Value, validators []Validator) error {
	fieldKind := field.Index(0).Kind()

	switch fieldKind { // nolint:exhaustive
	case reflect.Int:
		values := field.Interface().([]int)

		for _, value := range values {
			if err := validateInt(int64(value), validators); err != nil {
				return err
			}
		}
	case reflect.String:
		values := field.Interface().([]string)

		for _, value := range values {
			if err := validateString(value, validators); err != nil {
				return err
			}
		}
	default:
		return ErrNotImplementedType
	}

	return nil
}

func validateString(value string, validators []Validator) error {
	for _, validator := range validators {
		switch validator.key {
		case "regexp":
			match, err := regexp.MatchString(validator.value, value)
			if err != nil {
				return err
			}

			if !match {
				return ErrStringRegexp
			}
		case "len":
			length, err := strconv.Atoi(validator.value)
			if err != nil {
				return err
			}

			if len([]rune(value)) != length {
				return ErrStringLen
			}
		case "in":
			validatorValues := regexpParseValidatorValue.FindAllString(validator.value, -1)

			for _, word := range validatorValues {
				if word == value {
					return nil
				}
			}

			return ErrStringIn
		default:
			return ErrNotImplementedTag
		}
	}

	return nil
}

func validateInt(value int64, validators []Validator) error {
	for _, validator := range validators {
		switch validator.key {
		case "min":
			min, err := strconv.Atoi(validator.value)
			if err != nil {
				return err
			}

			if value < int64(min) {
				return ErrIntMin
			}

		case "max":
			max, err := strconv.Atoi(validator.value)
			if err != nil {
				return err
			}

			if value > int64(max) {
				return ErrIntMax
			}

		case "in":
			validatorValues := regexpParseValidatorValue.FindAllString(validator.value, -1)

			for _, strInt := range validatorValues {
				integer, err := strconv.Atoi(strInt)
				if err != nil {
					return err
				}

				if value == int64(integer) {
					return nil
				}
			}

			return ErrIntIn
		default:
			return ErrNotImplementedTag
		}
	}

	return nil
}
