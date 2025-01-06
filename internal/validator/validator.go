package validator

import "regexp"

var (
	// EmailRX is a regex for sanity checking the format of email addresses.
	// The regex pattern used is taken from  https://html.spec.whatwg.org/#valid-e-mail-address.
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// ImageNameRX is regex pattern to validate image name format
	// Format: filename-UUID_timestamp
	ImageNameRX = regexp.MustCompile(`^([a-z0-9-]+-\w+-\d{8}_\d{6})$`)

	// ImageFileNameRX is regex pattern to validate image file name format
	// almost the same as file name with extension postfix
	// Format: filename-UUID_timestamp.extension
	ImageFileNameRX = regexp.MustCompile(`^([a-z0-9-]+-\w+-\d{8}_\d{6})\.(jpeg|jpg|png|webp|gif)$`)
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = message
	}
}

func (v *Validator) ResetErrors() {
	v.Errors = make(map[string]string)
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}

	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
