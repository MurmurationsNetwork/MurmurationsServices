package core

// ProfileFetchError represents an error encountered when trying to fetch a profile.
type ProfileFetchError struct {
	Reason string
}

// Error function to represent the ProfileFetchError as a string.
func (e ProfileFetchError) Error() string {
	return "Profile fetch error: " + e.Reason
}
