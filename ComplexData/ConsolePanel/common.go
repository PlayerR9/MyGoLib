package ConsolePanel

// BoolFS returns "Yes" if val is true, "No" otherwise.
//
// Parameters:
//   - val: The boolean value.
//
// Returns:
//   - string: "Yes" if val is true, "No" otherwise.
//   - error: nil
func BoolFString(val bool) (string, error) {
	var res string

	if val {
		res = "Yes"
	} else {
		res = "No"
	}

	return res, nil
}
