package gdal

func (d SwigcptrDriver) CreateCopy(path string, ds Dataset, args ...interface{}) (ret Dataset, err Err) {
	ret = d.wrap_CreateCopy(path, ds, args...)
	if ret != nil {
		return
	}
	err = lastError()
	ret = nil
	return
}

func (d SwigcptrDriver) Create(path string, xsize, ysize, bands, etype int, options []string) (ret Dataset, err error) {
	ret = d.wrap_Create(path, xsize, ysize, bands, etype, options)
	if ret != nil {
		return
	}
	err = lastError()
	ret = nil
	return
}
