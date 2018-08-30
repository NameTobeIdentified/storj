package testsuite

import (
	"math/rand"
	"testing"

	"storj.io/storj/storage"
)

func testIterate(t *testing.T, store storage.KeyValueStore) {
	items := storage.Items{
		newItem("a", "a", false),
		newItem("b/1", "b/1", false),
		newItem("b/2", "b/2", false),
		newItem("b/3", "b/3", false),
		newItem("c", "c", false),
		newItem("c/", "c/", false),
		newItem("c//", "c//", false),
		newItem("c/1", "c/1", false),
		newItem("g", "g", false),
		newItem("h", "h", false),
	}
	rand.Shuffle(len(items), items.Swap)
	defer cleanupItems(store, items)

	for _, item := range items {
		if err := store.Put(item.Key, item.Value); err != nil {
			t.Fatalf("failed to put: %v", err)
		}
	}

	t.Run("no limits", func(t *testing.T) {
		store.Iterate(nil, nil, '/',
			checkIterator(t, storage.Items{
				newItem("a", "a", false),
				newItem("b/", "", true),
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at a", func(t *testing.T) {
		store.Iterate(nil, storage.Key("a"), '/',
			checkIterator(t, storage.Items{
				newItem("a", "a", false),
				newItem("b/", "", true),
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after a", func(t *testing.T) {
		store.Iterate(nil, storage.NextKey(storage.Key("a")), '/',
			checkIterator(t, storage.Items{
				newItem("b/", "", true),
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at b", func(t *testing.T) {
		store.Iterate(nil, storage.Key("b"), '/',
			checkIterator(t, storage.Items{
				newItem("b/", "", true),
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after b", func(t *testing.T) {
		store.Iterate(nil, storage.NextKey(storage.Key("b")), '/',
			checkIterator(t, storage.Items{
				newItem("b/", "", true),
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at c", func(t *testing.T) {
		store.Iterate(nil, storage.Key("c"), '/',
			checkIterator(t, storage.Items{
				newItem("c", "c", false),
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after c", func(t *testing.T) {
		store.Iterate(nil, storage.NextKey(storage.Key("c")), '/',
			checkIterator(t, storage.Items{
				newItem("c/", "", true),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at e", func(t *testing.T) {
		store.Iterate(nil, storage.Key("e"), '/',
			checkIterator(t, storage.Items{
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after e", func(t *testing.T) {
		store.Iterate(nil, storage.NextKey(storage.Key("e")), '/',
			checkIterator(t, storage.Items{
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("prefix b slash", func(t *testing.T) {
		store.Iterate(storage.Key("b/"), nil, '/',
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
			}))
	})

	t.Run("prefix c slash", func(t *testing.T) {
		store.Iterate(storage.Key("c/"), nil, '/',
			checkIterator(t, storage.Items{
				newItem("c/", "c/", false),
				newItem("c//", "", true),
				newItem("c/1", "c/1", false),
			}))
	})

	t.Run("prefix c slash slash", func(t *testing.T) {
		store.Iterate(storage.Key("c//"), nil, '/',
			checkIterator(t, storage.Items{
				newItem("c//", "c//", false),
			}))
	})
}

func testIterateAll(t *testing.T, store storage.KeyValueStore) {
	items := storage.Items{
		newItem("a", "a", false),
		newItem("b/1", "b/1", false),
		newItem("b/2", "b/2", false),
		newItem("b/3", "b/3", false),
		newItem("c", "c", false),
		newItem("c/", "c/", false),
		newItem("c//", "c//", false),
		newItem("c/1", "c/1", false),
		newItem("g", "g", false),
		newItem("h", "h", false),
	}
	rand.Shuffle(len(items), items.Swap)
	defer cleanupItems(store, items)

	for _, item := range items {
		if err := store.Put(item.Key, item.Value); err != nil {
			t.Fatalf("failed to put: %v", err)
		}
	}

	t.Run("no limits", func(t *testing.T) {
		store.IterateAll(nil, nil,
			checkIterator(t, storage.Items{
				newItem("a", "a", false),
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at a", func(t *testing.T) {
		store.IterateAll(nil, storage.Key("a"),
			checkIterator(t, storage.Items{
				newItem("a", "a", false),
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after a", func(t *testing.T) {
		store.IterateAll(nil, storage.NextKey(storage.Key("a")),
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at b", func(t *testing.T) {
		store.IterateAll(nil, storage.Key("b"),
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after b", func(t *testing.T) {
		store.IterateAll(nil, storage.NextKey(storage.Key("b")),
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at c", func(t *testing.T) {
		store.IterateAll(nil, storage.Key("c"),
			checkIterator(t, storage.Items{
				newItem("c", "c", false),
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("after c", func(t *testing.T) {
		store.IterateAll(nil, storage.NextKey(storage.Key("c")),
			checkIterator(t, storage.Items{
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("at e", func(t *testing.T) {
		store.IterateAll(nil, storage.Key("e"),
			checkIterator(t, storage.Items{
				newItem("g", "g", false),
				newItem("h", "h", false),
			}))
	})

	t.Run("prefix b slash", func(t *testing.T) {
		store.IterateAll(storage.Key("b/"), nil,
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
			}))

		store.IterateAll(storage.Key("b/"), storage.Key("a"),
			checkIterator(t, storage.Items{
				newItem("b/1", "b/1", false),
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
			}))

		store.IterateAll(storage.Key("b/"), storage.Key("b/2"),
			checkIterator(t, storage.Items{
				newItem("b/2", "b/2", false),
				newItem("b/3", "b/3", false),
			}))
	})

	t.Run("prefix c slash", func(t *testing.T) {
		store.IterateAll(storage.Key("c/"), nil,
			checkIterator(t, storage.Items{
				newItem("c/", "c/", false),
				newItem("c//", "c//", false),
				newItem("c/1", "c/1", false),
			}))
	})

	t.Run("prefix c slash slash", func(t *testing.T) {
		store.IterateAll(storage.Key("c//"), nil,
			checkIterator(t, storage.Items{
				newItem("c//", "c//", false),
			}))
	})
}
