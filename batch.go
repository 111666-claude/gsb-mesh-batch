// Package batch 是网格合批：按材质分组、顶点上限拆批与顺序稳定。
// 缺陷：不按材质分组、不看顶点上限、空网格不跳过、批次顺序不稳定、
// 每次加入都要重扫已有批次。
package batch

// Mesh 是一个待合批的网格。
type Mesh struct {
	ID       string
	Material string
	Verts    int
}

// Batcher 是合批器。
type Batcher struct {
	maxVerts int
	meshes   []Mesh
	scanned  int
}

// New 建合批器。
func New(maxVerts int) *Batcher { return &Batcher{maxVerts: maxVerts} }

// Add 加入一个网格。
func (b *Batcher) Add(m Mesh) {
	b.scanned += len(b.meshes)
	b.meshes = append(b.meshes, m)
}

// Batches 返回每批的网格 id 列表。
func (b *Batcher) Batches() [][]string {
	b.scanned += len(b.meshes)
	out := [][]string{}
	if len(b.meshes) == 0 {
		return out
	}
	ids := []string{}
	for _, m := range b.meshes {
		ids = append(ids, m.ID)
	}
	return append(out, ids)
}

// Scanned 是累计扫描量（规模观测）。
func (b *Batcher) Scanned() int { return b.scanned }
