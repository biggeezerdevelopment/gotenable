package base

import (
	"context"
	"encoding/json"
	"fmt"
)

// PaginationInfo contains pagination metadata from API responses.
type PaginationInfo struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PageFetcher is a function that fetches a page of results.
type PageFetcher func(ctx context.Context, offset, limit int) (json.RawMessage, *PaginationInfo, error)

// ItemTransformer is a function that transforms raw JSON into typed items.
type ItemTransformer[T any] func(data json.RawMessage) ([]T, error)

// Iterator provides pagination over API results.
type Iterator[T any] struct {
	ctx         context.Context
	fetcher     PageFetcher
	transformer ItemTransformer[T]
	limit       int
	offset      int
	total       int
	page        []T
	pageIndex   int
	count       int
	maxPages    int
	pagesLoaded int
	done        bool
	err         error
	current     T
}

// IteratorOption configures an Iterator.
type IteratorOption[T any] func(*Iterator[T])

// NewIterator creates a new paginated iterator.
func NewIterator[T any](
	ctx context.Context,
	fetcher PageFetcher,
	transformer ItemTransformer[T],
	opts ...IteratorOption[T],
) *Iterator[T] {
	it := &Iterator[T]{
		ctx:         ctx,
		fetcher:     fetcher,
		transformer: transformer,
		limit:       100,
		offset:      0,
		total:       -1, // Unknown until first fetch
		maxPages:    0,  // 0 means no limit
	}

	for _, opt := range opts {
		opt(it)
	}

	return it
}

// WithLimit sets the page size limit.
func WithLimit[T any](limit int) IteratorOption[T] {
	return func(it *Iterator[T]) {
		it.limit = limit
	}
}

// WithOffset sets the starting offset.
func WithOffset[T any](offset int) IteratorOption[T] {
	return func(it *Iterator[T]) {
		it.offset = offset
	}
}

// WithMaxPages sets the maximum number of pages to fetch.
func WithMaxPages[T any](maxPages int) IteratorOption[T] {
	return func(it *Iterator[T]) {
		it.maxPages = maxPages
	}
}

// Next returns the next item. Returns false when iteration is complete.
func (it *Iterator[T]) Next() bool {
	if it.done || it.err != nil {
		return false
	}

	// Check if we've reached the total
	if it.total >= 0 && it.count >= it.total {
		it.done = true
		return false
	}

	// Check if we need to fetch more data
	if it.pageIndex >= len(it.page) {
		if !it.fetchNextPage() {
			return false
		}
		return true // fetchNextPage already set the current item
	}

	// Get current item and advance
	it.current = it.page[it.pageIndex]
	it.pageIndex++
	it.count++
	return true
}

// Item returns the current item.
func (it *Iterator[T]) Item() T {
	return it.current
}

// Err returns any error that occurred during iteration.
func (it *Iterator[T]) Err() error {
	return it.err
}

// Count returns the number of items returned so far.
func (it *Iterator[T]) Count() int {
	return it.count
}

// Total returns the total number of items available (-1 if unknown).
func (it *Iterator[T]) Total() int {
	return it.total
}

// fetchNextPage fetches the next page of results.
func (it *Iterator[T]) fetchNextPage() bool {
	// Check max pages limit
	if it.maxPages > 0 && it.pagesLoaded >= it.maxPages {
		it.done = true
		return false
	}

	// Check if we've already fetched everything
	if it.total >= 0 && it.offset >= it.total {
		it.done = true
		return false
	}

	// Fetch the page
	data, pagination, err := it.fetcher(it.ctx, it.offset, it.limit)
	if err != nil {
		it.err = err
		return false
	}

	// Update pagination info
	if pagination != nil {
		it.total = pagination.Total
	}

	// Transform the data
	items, err := it.transformer(data)
	if err != nil {
		it.err = fmt.Errorf("failed to transform page data: %w", err)
		return false
	}

	if len(items) == 0 {
		it.done = true
		return false
	}

	it.page = items
	it.pageIndex = 0
	it.offset += len(items)
	it.pagesLoaded++

	// Get first item
	it.current = it.page[it.pageIndex]
	it.pageIndex++
	it.count++
	return true
}

// All returns all remaining items as a slice.
func (it *Iterator[T]) All() ([]T, error) {
	var results []T
	for it.Next() {
		results = append(results, it.Item())
	}
	if it.err != nil {
		return nil, it.err
	}
	return results, nil
}

// Take returns up to n items.
func (it *Iterator[T]) Take(n int) ([]T, error) {
	var results []T
	for i := 0; i < n && it.Next(); i++ {
		results = append(results, it.Item())
	}
	if it.err != nil {
		return nil, it.err
	}
	return results, nil
}

// ForEach calls the given function for each item.
func (it *Iterator[T]) ForEach(fn func(T) error) error {
	for it.Next() {
		if err := fn(it.Item()); err != nil {
			return err
		}
	}
	return it.err
}

// Channel returns a channel that yields items.
func (it *Iterator[T]) Channel() <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for it.Next() {
			select {
			case ch <- it.Item():
			case <-it.ctx.Done():
				return
			}
		}
	}()
	return ch
}
