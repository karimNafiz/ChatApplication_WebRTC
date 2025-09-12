package validator

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}

func (v *Validator) Unique(values []string) bool {
	set := make(map[string]bool)
	for i := range values {
		set[values[i]] = true
	}

	return len(set) == len(values)
}
