package types

import (
	"github.com/dlclark/regexp2"
)

type Regexp struct {
	*regexp2.Regexp
}

func RegexpFromString(pattern string) (*Regexp, error) {
	regexp, err := regexp2.Compile(pattern, regexp2.DefaultUnmarshalOptions)
	if err != nil {
		return nil, err
	}
	return &Regexp{
		regexp,
	}, nil
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := regexp2.Compile(string(b), regexp2.DefaultUnmarshalOptions)
	if err != nil {
		return err
	}

	r.Regexp = regex

	return nil
}

// MarshalText marshals regexp.Regexp as string
func (r *Regexp) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
}
