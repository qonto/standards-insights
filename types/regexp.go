package types

import "regexp"

type Regexp struct {
	*regexp.Regexp
}

func RegexpFromString(pattern string) (*Regexp, error) {
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Regexp{
		regexp,
	}, nil
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := regexp.Compile(string(b))
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
