package base

type ModelPopulater[M any] interface {
	PopulateModel(*M)
}

type ModelRefPopulater[M any] interface {
	PopulateModelRef(*M)
}
