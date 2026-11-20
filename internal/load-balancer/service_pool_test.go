package loadbalancer

import (
	"slices"
	"testing"
)

func TestServicePoolReplaceUpsertRemove(t *testing.T) {
	pool := NewServicePool()

	pool.Replace([]ServiceInstance{{ID: "a", Weight: 1}, {ID: "b", Weight: 1}})
	pool.Upsert(ServiceInstance{ID: "b", Weight: 3})
	pool.Upsert(ServiceInstance{ID: "c", Weight: 1})
	pool.Remove("a")
	pool.Remove("unknown")

	want := []ServiceInstance{{ID: "b", Weight: 3}, {ID: "c", Weight: 1}}
	if got := pool.Instances(); !slices.Equal(got, want) {
		t.Errorf("Instances() = %v, want %v", got, want)
	}
}

func TestServicePoolReturnedInstancesAreNotModifiedLater(t *testing.T) {
	input := []ServiceInstance{{ID: "a", Weight: 1}, {ID: "b", Weight: 1}}

	pool := NewServicePool()
	pool.Replace(input)
	input[0].Weight = 99

	snapshot := pool.Instances()
	want := []ServiceInstance{{ID: "a", Weight: 1}, {ID: "b", Weight: 1}}

	pool.Upsert(ServiceInstance{ID: "a", Weight: 5})
	pool.Remove("b")

	if !slices.Equal(snapshot, want) {
		t.Errorf("earlier Instances() result = %v, want it unchanged %v", snapshot, want)
	}
}
