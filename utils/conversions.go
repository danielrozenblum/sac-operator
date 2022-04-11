package utils

func Convert_Pointer_bool_To_bool_with_default(in *bool, defaultOutput bool) bool {
	if in == nil {
		return defaultOutput
	}
	return *in
}

func Convert_bool_To_Pointer_bool(in bool) *bool {
	return &in
}
