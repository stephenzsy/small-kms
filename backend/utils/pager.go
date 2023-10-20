package utils

import (
	"context"
)

type Pager[T any] interface {
	More() bool
	NextPage(c context.Context) (page T, err error)
}

type ItemsPager[T any] Pager[[]T]

type mappedPager[T any, U any] struct {
	from    Pager[U]
	mapFunc MapFunc[T, U]
}

func NewMappedPager[T any, U any](from Pager[U], mapFunc MapFunc[T, U]) Pager[T] {
	return &mappedPager[T, U]{from: from, mapFunc: mapFunc}
}

func (p *mappedPager[T, U]) More() bool {
	return p.from.More()
}

func (p *mappedPager[T, U]) NextPage(c context.Context) (items T, err error) {
	r, err := p.from.NextPage(c)
	return p.mapFunc(r), err
}

type mappedItemsPager[T any, U any] struct {
	from    ItemsPager[U]
	mapFunc MapFunc[T, U]
}

// More implements Pager.
func (p *mappedItemsPager[T, U]) More() bool {
	return p.from.More()
}

// NextPage implements Pager.
func (p *mappedItemsPager[T, U]) NextPage(c context.Context) (items []T, err error) {
	var fromSlice []U
	if fromSlice, err = p.from.NextPage(c); err != nil {
		return
	}
	items = MapSlice(fromSlice, p.mapFunc)
	return
}

var _ ItemsPager[string] = (*mappedItemsPager[string, any])(nil)

func NewMappedItemsPager[T any, U any](from ItemsPager[U], f MapFunc[T, U]) ItemsPager[T] {
	return &mappedItemsPager[T, U]{from: from, mapFunc: f}
}

func PagerAllItems[T any](pager ItemsPager[T], ctx context.Context) (items []T, err error) {
	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, t...)
	}
	return
}

func ReservedFirst[T any](
	allItems []T,
	reservedDefaults []T,
	mapReservedIndex func(T) int) (items []T) {
	items = make([]T, len(reservedDefaults), len(allItems)+len(reservedDefaults))
	copy(items, reservedDefaults)
	for _, item := range allItems {
		if mappedIndex := mapReservedIndex(item); mappedIndex >= 0 && mappedIndex < len(items) {
			items[mappedIndex] = item
			continue
		} else {
			items = append(items, item)
		}
	}
	return
}
