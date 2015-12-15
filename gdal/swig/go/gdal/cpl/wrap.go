package cpl

func ReadDir(utf8_path string) (ret []string, err error) {
	ret = wrap_ReadDir(utf8_path)
	err = LastError()
	return
}
