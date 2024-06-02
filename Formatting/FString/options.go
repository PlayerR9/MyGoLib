package FString

// Settinger is an interface that represents the settings for the formatting
// functions.
type Settinger interface{}

// Option is a function that sets the settings for the formatting functions.
//
// Parameters:
//   - Settinger: The settings to set.
type Option func(Settinger)
