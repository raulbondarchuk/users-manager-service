package utils

// StringOrNil returns nil if the input string is empty, otherwise returns a pointer to the string
func StringOrNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
