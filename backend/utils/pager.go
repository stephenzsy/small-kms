package utils

import (
	"context"
	"encoding/json"
)

type Pager[T any] interface {
	More() bool
	NextPage(c context.Context) (page T, err error)
}

type ItemsPager[T any] interface {
	Pager[[]T]
}

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
	return PagerToSlice(ctx, pager)
}

func PagerToSlice[T any](ctx context.Context, pager ItemsPager[T]) (items []T, err error) {
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

type chainedItemPagers[T any] struct {
	pagers []ItemsPager[T]
	index  int
}

// More implements ItemsPager.
func (p *chainedItemPagers[T]) More() bool {
	return p.index < len(p.pagers)
}

// NextPage implements ItemsPager.
func (p *chainedItemPagers[T]) NextPage(c context.Context) (page []T, err error) {
	for p.index < len(p.pagers) && !p.pagers[p.index].More() {
		p.index++
	}
	if p.index >= len(p.pagers) {
		return nil, nil
	}
	return p.pagers[p.index].NextPage(c)
}

func NewChainedItemPagers[T any](pagers ...ItemsPager[T]) ItemsPager[T] {
	return &chainedItemPagers[T]{pagers: pagers}
}

type SerializableItemsPager[T any] struct {
	ItemsPager[T]
	ctx context.Context
}

// MarshalJSON implements json.Marshaler.
func (p *SerializableItemsPager[T]) MarshalJSON() ([]byte, error) {
	b := append([]byte(nil), '[')
	for p.More() {
		items, err := p.NextPage(p.ctx)
		if err != nil {
			return nil, err
		}
		for i, item := range items {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			if i > 0 {
				b = append(b, ',')
			}
			b = append(b, itemBytes...)
		}
	}
	b = append(b, ']')
	return b, nil
}

var _ json.Marshaler = (*SerializableItemsPager[any])(nil)

func NewSerializableItemsPager[T any](ctx context.Context, pager ItemsPager[T]) *SerializableItemsPager[T] {
	return &SerializableItemsPager[T]{ctx: ctx, ItemsPager: pager}
}
