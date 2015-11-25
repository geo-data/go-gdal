package gdal

func (d SwigcptrDriver) CreateCopy(filename string, ds Dataset, args ...interface{}) (ret Dataset, err Err) {
	ret = d.Wrap_CreateCopy(filename, ds, args...)
	if ret != nil {
		return
	}
	err = lastError()
	ret = nil
	return
}
