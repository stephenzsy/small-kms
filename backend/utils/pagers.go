package utils

import (
	"context"
	"encoding/json"
)

type Pager[T any] interface {
	More() bool
	NextPage() (page T, err error)
}

type PagerRequireContext[T any] interface {
	More() bool
	NextPage(context.Context) (page T, err error)
}

type PagerWithContext[T any] struct {
	pager PagerRequireContext[T]
	ctx   context.Context
}

func (p *PagerWithContext[T]) More() bool {
	return p.pager.More()
}

func (p *PagerWithContext[T]) NextPage() (page T, err error) {
	return p.pager.NextPage(p.ctx)
}

func NewPagerWithContext[T any](pager PagerRequireContext[T], ctx context.Context) Pager[T] {
	return &PagerWithContext[T]{pager: pager, ctx: ctx}
}

type ItemsPager[T any] interface {
	Pager[[]T]
}

// map a pager to another
type mappedPager[T, U any] struct {
	from    Pager[U]
	mapFunc MapFunc[T, U]
}

func NewMappedPager[T any, U any](from Pager[U], mapFunc MapFunc[T, U]) Pager[T] {
	return &mappedPager[T, U]{from: from, mapFunc: mapFunc}
}

func (p *mappedPager[T, U]) More() bool {
	return p.from.More()
}

func (p *mappedPager[T, U]) NextPage() (items T, err error) {
	r, err := p.from.NextPage()
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
func (p *mappedItemsPager[T, U]) NextPage() (items []T, err error) {
	var fromSlice []U
	if fromSlice, err = p.from.NextPage(); err != nil {
		return
	}
	items = MapSlice(fromSlice, p.mapFunc)
	return
}

var _ ItemsPager[string] = (*mappedItemsPager[string, any])(nil)

func NewMappedItemsPager[T, U any](from ItemsPager[U], f MapFunc[T, U]) ItemsPager[T] {
	return &mappedItemsPager[T, U]{from: from, mapFunc: f}
}

func PagerToSlice[T any](pager ItemsPager[T]) (items []T, err error) {
	for pager.More() {
		t, scanErr := pager.NextPage()
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, t...)
	}
	return
}

type chainedItemPagers[T any] struct {
	chainedPagers []ItemsPager[T]
	currentIndex  int
}

// More implements ItemsPager.
func (p *chainedItemPagers[T]) More() bool {
	return p.currentIndex < len(p.chainedPagers)
}

// NextPage implements ItemsPager.
func (p *chainedItemPagers[T]) NextPage() (page []T, err error) {
	for p.currentIndex < len(p.chainedPagers) && !p.chainedPagers[p.currentIndex].More() {
		p.currentIndex++
	}
	if p.currentIndex >= len(p.chainedPagers) {
		return nil, nil
	}
	return p.chainedPagers[p.currentIndex].NextPage()
}

func NewChainedItemPagers[T any](pagers ...ItemsPager[T]) ItemsPager[T] {
	return &chainedItemPagers[T]{chainedPagers: pagers}
}

type SerializableItemsPager[T any] struct {
	ItemsPager[T]
}

// MarshalJSON implements json.Marshaler.
func (p *SerializableItemsPager[T]) MarshalJSON() ([]byte, error) {
	b := append([]byte(nil), '[')
	for p.More() {
		items, err := p.NextPage()
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

func NewSerializableItemsPager[T any](pager ItemsPager[T]) *SerializableItemsPager[T] {
	return &SerializableItemsPager[T]{ItemsPager: pager}
}
