package base

import (
	"context"
	"encoding/json"
	"testing"
)

type testItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestIterator(t *testing.T) {
	// Simulate paginated data
	allItems := []testItem{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
		{ID: 3, Name: "three"},
		{ID: 4, Name: "four"},
		{ID: 5, Name: "five"},
	}

	createFetcher := func() PageFetcher {
		return func(ctx context.Context, offset, limit int) (json.RawMessage, *PaginationInfo, error) {
			end := offset + limit
			if end > len(allItems) {
				end = len(allItems)
			}
			if offset >= len(allItems) {
				data, _ := json.Marshal([]testItem{})
				return data, &PaginationInfo{Total: len(allItems), Limit: limit, Offset: offset}, nil
			}

			items := allItems[offset:end]
			data, _ := json.Marshal(items)
			return data, &PaginationInfo{
				Total:  len(allItems),
				Limit:  limit,
				Offset: offset,
			}, nil
		}
	}

	transformer := func(data json.RawMessage) ([]testItem, error) {
		var items []testItem
		err := json.Unmarshal(data, &items)
		return items, err
	}

	ctx := context.Background()

	t.Run("iterate all items", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer, WithLimit[testItem](2))

		var collected []testItem
		for iter.Next() {
			collected = append(collected, iter.Item())
		}

		if err := iter.Err(); err != nil {
			t.Fatalf("Iterator error: %v", err)
		}

		if len(collected) != len(allItems) {
			t.Errorf("Expected %d items, got %d", len(allItems), len(collected))
		}

		for i, item := range collected {
			if item.ID != allItems[i].ID {
				t.Errorf("Item %d: expected ID %d, got %d", i, allItems[i].ID, item.ID)
			}
		}
	})

	t.Run("Take() method", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer, WithLimit[testItem](2))

		collected, err := iter.Take(3)
		if err != nil {
			t.Fatalf("Take() error: %v", err)
		}

		if len(collected) != 3 {
			t.Errorf("Expected 3 items, got %d", len(collected))
		}
	})

	t.Run("All() method", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer, WithLimit[testItem](2))

		collected, err := iter.All()
		if err != nil {
			t.Fatalf("All() error: %v", err)
		}

		if len(collected) != len(allItems) {
			t.Errorf("Expected %d items, got %d", len(allItems), len(collected))
		}
	})

	t.Run("ForEach() method", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer, WithLimit[testItem](2))

		count := 0
		err := iter.ForEach(func(item testItem) error {
			count++
			return nil
		})

		if err != nil {
			t.Fatalf("ForEach() error: %v", err)
		}

		if count != len(allItems) {
			t.Errorf("Expected %d iterations, got %d", len(allItems), count)
		}
	})

	t.Run("max pages limit", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer,
			WithLimit[testItem](2),
			WithMaxPages[testItem](1),
		)

		var collected []testItem
		for iter.Next() {
			collected = append(collected, iter.Item())
		}

		if err := iter.Err(); err != nil {
			t.Fatalf("Iterator error: %v", err)
		}

		// With limit 2 and max pages 1, should get only 2 items
		if len(collected) != 2 {
			t.Errorf("Expected 2 items with max pages 1, got %d", len(collected))
		}
	})

	t.Run("offset option", func(t *testing.T) {
		iter := NewIterator(ctx, createFetcher(), transformer,
			WithLimit[testItem](2),
			WithOffset[testItem](2),
		)

		collected, err := iter.All()
		if err != nil {
			t.Fatalf("All() error: %v", err)
		}

		// Starting at offset 2, should get items 3, 4, 5
		if len(collected) != 3 {
			t.Errorf("Expected 3 items starting at offset 2, got %d", len(collected))
		}

		if len(collected) > 0 && collected[0].ID != 3 {
			t.Errorf("First item should have ID 3, got %d", collected[0].ID)
		}
	})
}

func TestIteratorCount(t *testing.T) {
	allItems := []testItem{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
		{ID: 3, Name: "three"},
	}

	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *PaginationInfo, error) {
		if offset >= len(allItems) {
			data, _ := json.Marshal([]testItem{})
			return data, &PaginationInfo{Total: len(allItems)}, nil
		}
		end := offset + limit
		if end > len(allItems) {
			end = len(allItems)
		}
		items := allItems[offset:end]
		data, _ := json.Marshal(items)
		return data, &PaginationInfo{Total: len(allItems), Limit: limit, Offset: offset}, nil
	}

	transformer := func(data json.RawMessage) ([]testItem, error) {
		var items []testItem
		err := json.Unmarshal(data, &items)
		return items, err
	}

	ctx := context.Background()
	iter := NewIterator(ctx, fetcher, transformer, WithLimit[testItem](10))

	// Count should start at 0
	if iter.Count() != 0 {
		t.Errorf("Initial count should be 0, got %d", iter.Count())
	}

	// After first Next()
	iter.Next()
	if iter.Count() != 1 {
		t.Errorf("Count after first Next() should be 1, got %d", iter.Count())
	}

	// After iterating all
	for iter.Next() {
	}
	if iter.Count() != len(allItems) {
		t.Errorf("Final count should be %d, got %d", len(allItems), iter.Count())
	}
}

func TestIteratorEmpty(t *testing.T) {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *PaginationInfo, error) {
		data, _ := json.Marshal([]testItem{})
		return data, &PaginationInfo{Total: 0}, nil
	}

	transformer := func(data json.RawMessage) ([]testItem, error) {
		var items []testItem
		err := json.Unmarshal(data, &items)
		return items, err
	}

	ctx := context.Background()
	iter := NewIterator(ctx, fetcher, transformer)

	// Should return false immediately for empty results
	if iter.Next() {
		t.Error("Next() should return false for empty iterator")
	}

	if iter.Count() != 0 {
		t.Errorf("Count should be 0, got %d", iter.Count())
	}
}
