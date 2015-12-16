package cpl

func ReadDir(utf8_path string) (ret []string, err error) {
	defer ErrorTrap()(&err)
	ret = wrap_ReadDir(utf8_path)
	return
}
