package helpers

func Int64Pointer(v int64) *int64 {
	return &v
}

func IntPointer(v int) *int {
	return &v
}

func StringPointer(v string) *string {
	if v == "" {
		return nil
	}

	return &v
}
