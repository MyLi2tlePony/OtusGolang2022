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

	ErrItemIsNotPointer = errors.New("item must be a pointer")
	ErrStringRegexp     = errors.New("error string regexp")
	ErrStringLen        = errors.New("error string len")
	ErrStringIn         = errors.New("error string in")
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
	value := reflect.ValueOf(item).Elem()
	if !value.CanAddr() {
		return ErrItemIsNotPointer
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
	if kind == reflect.String {
		value := field.String()
		for _, validator := range validators {
			err := validateString(value, validator)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateString(value string, validator Validator) error {
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
		b, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if len([]rune(value)) != b {
			return ErrStringLen
		}
	case "in":
		validatorValues := regexpParseValidatorValue.FindAllString(value, -1)

		for _, word := range validatorValues {
			if word == value {
				return nil
			}
		}

		return ErrStringIn
	}

	return nil
}
