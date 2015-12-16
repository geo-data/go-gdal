package cpl

func getLastError(err *error) (set bool) {
	e := LastError()
	if e == nil {
		return
	}

	*err = error(e)
	set = true
	return
}

func ErrorTrap() func(err *error) (set bool) {
	ErrorReset()
	return getLastError
}
