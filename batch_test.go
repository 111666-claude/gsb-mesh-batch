package batch

import "testing"

func TestEmptyBatcher(t *testing.T) {
	if len(New(100).Batches()) != 0 {
		t.Fatal("没有网格时不该有批次")
	}
}

func TestSingleMeshIsOneBatch(t *testing.T) {
	batcher := New(100)
	batcher.Add(Mesh{ID: "a", Material: "stone", Verts: 10})
	batches := batcher.Batches()
	if len(batches) != 1 || len(batches[0]) != 1 {
		t.Fatal("单个网格应该是一个批次")
	}
}

func TestMaxVertsIsStored(t *testing.T) {
	if New(64).maxVerts != 64 {
		t.Fatal("顶点上限应该能读回来")
	}
}

func TestSetLimitIsStored(t *testing.T) {
	batcher := New(64)
	batcher.SetLimit(8)
	if batcher.maxVerts != 8 {
		t.Fatal("改上限之后应该能读回来")
	}
}

func TestRemoveUnknownReturnsFalse(t *testing.T) {
	batcher := New(100)
	batcher.Add(Mesh{ID: "a", Material: "stone", Verts: 10})
	if batcher.Remove("nope") {
		t.Fatal("移除没加入过的网格应该返回 false")
	}
}

func TestBatchesResultIsCopy(t *testing.T) {
	batcher := New(100)
	batcher.Add(Mesh{ID: "a", Material: "stone", Verts: 10})
	batches := batcher.Batches()
	batches[0][0] = "tampered"
	if got := batcher.Batches()[0][0]; got != "a" {
		t.Fatalf("改动返回结果不该影响内部批次表，得到 %s", got)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	batcher := New(64)
	if batcher.Scanned() != 0 || batcher.Recomputed() != 0 {
		t.Fatal("两个计数器初始都应该是 0")
	}
}
