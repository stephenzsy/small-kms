package utils

import (
	"context"
)

type ItemsPager[T any] interface {
	More() bool
	NextPage(c context.Context) (items []T, err error)
}

type wrappedItemsPager[T any] struct {
	items   []T
	hasRead bool
}

func (p *wrappedItemsPager[T]) More() bool {
	return !p.hasRead
}

func (p *wrappedItemsPager[T]) NextPage(c context.Context) (items []T, err error) {
	p.hasRead = true
	return p.items, nil
}

func NewWrappedPager[T any](items []T) ItemsPager[T] {
	return &wrappedItemsPager[T]{items: items}
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
	items = MapSlices(fromSlice, p.mapFunc)
	return
}

var _ ItemsPager[string] = (*mappedItemsPager[string, any])(nil)

func NewMappedPager[T any, U any](from ItemsPager[U], f MapFunc[T, U]) ItemsPager[T] {
	return &mappedItemsPager[T, U]{from: from, mapFunc: f}
}

func PagerAllItems[T any](pager ItemsPager[T], ctx context.Context) (items []T, err error) {
	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			return nil, err
		}
		items = append(items, t...)
	}
	return
}
