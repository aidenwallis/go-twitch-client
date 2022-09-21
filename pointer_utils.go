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

// Pointer returns the passed value as a pointer. This is particularly useful in cases where you need to
// pass a pointer to a given func input, for example:
//
//	value := "string"
//	valuePointer := &value
//	log.Println(twitch.PointerValue(valuePointer)) // returns "string"
//
//	valuePointer = nil // set it to nil
//	log.Println(twitch.PointerValue(valuePointer)) // returns ""
func PointerValue[T any](v *T) T {
	if v == nil {
		var out T
		return out
	}
	return *v
}
