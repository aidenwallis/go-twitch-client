package twitch

// Pointer returns the passed value as a pointer. This is particularly useful in cases where you need to
// pass a pointer to a given func input, for example:
//
//	&helix.ModifyChannelInformation{
//	    Title: twitch.Pointer("my new title"),
//	}
func Pointer[T any](v T) *T {
	return &v
}
