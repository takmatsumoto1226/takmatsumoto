package ftn

type Filter interface {
	filter(*FTNBT)
	setNext(Filter)
}
