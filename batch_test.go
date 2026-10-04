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

func TestScannedStartsAtZero(t *testing.T) {
	if New(64).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
