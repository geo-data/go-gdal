package gdal

type NameProcessor interface {
	DestName(name string) (ds Dataset, err error)
	SetDatasets(datasets []Dataset)
	Datasets() []Dataset
}

type DatasetProcessor interface {
	NameProcessor
	DestDS(output Dataset) (err error)
}
