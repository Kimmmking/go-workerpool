package workerpool

import "testing"

func TestNewPool(t *testing.T) {
	p := New(5)
	if p.capacity != 5 {
		t.Errorf("want 5, actual %d\n", p.capacity)
	}

	if !p.block {
		t.Errorf("want true, actual %t\n", p.block)
	}

	if p.preAlloc {
		t.Errorf("want false, actual %t\n", p.preAlloc)
	}

	if len(p.active) != 0 {
		t.Errorf("want 0, actual %d\n", len(p.active))
	}

	p.Free()
	if len(p.active) != 0 {
		t.Errorf("want 0, actual %d\n", len(p.active))
	}

	p = New(5, WithBlock(false), WithPreAllocWorkers(true))
	if p.block {
		t.Errorf("want false, actual %t\n", p.block)
	}

	if !p.preAlloc {
		t.Errorf("want true, actual %t\n", p.preAlloc)
	}
	if len(p.active) != 5 {
		t.Errorf("want 5, actual %d\n", len(p.active))
	}

	p.Free()
	if len(p.active) != 0 {
		t.Errorf("want 0, actual %d\n", len(p.active))
	}

	p = New(-1)
	if p.capacity != defaultCapacity {
		t.Errorf("want %d, actual %d\n", defaultCapacity, p.capacity)
	}
}

func TestSchedule(t *testing.T) {

}
