// Package batch 是网格合批：按材质分组、顶点上限拆批、加入顺序稳定、移除与改上限后重新紧凑化。
// 缺陷：所有网格塞进同一批、不看顶点上限、空材质照收、移除只认最后一个且不更新批次表、
// 改上限不重排、加入要重扫全部网格。
package batch

// Mesh 是一个待合批的网格。
type Mesh struct {
	ID       string
	Material string
	Verts    int
}

// Batcher 是合批器。
type Batcher struct {
	maxVerts   int
	meshes     []Mesh
	table      [][]string
	recomputed int
	scanned    int
}

// New 建合批器。
func New(maxVerts int) *Batcher { return &Batcher{maxVerts: maxVerts} }

// Add 加入一个网格。
func (b *Batcher) Add(m Mesh) {
	b.scanned += len(b.meshes)
	b.meshes = append(b.meshes, m)
	if len(b.table) == 0 {
		b.table = [][]string{{}}
	}
	b.table[0] = append(b.table[0], m.ID)
}

// Remove 移除一个网格，返回它原本是否存在。
func (b *Batcher) Remove(id string) bool {
	if len(b.meshes) == 0 || b.meshes[len(b.meshes)-1].ID != id {
		return false
	}
	b.meshes = b.meshes[:len(b.meshes)-1]
	return true
}

// SetLimit 改顶点上限。
func (b *Batcher) SetLimit(limit int) { b.maxVerts = limit }

// Batches 返回每批的网格 id 列表。
func (b *Batcher) Batches() [][]string {
	out := [][]string{}
	for _, item := range b.table {
		out = append(out, append([]string{}, item...))
	}
	return out
}

// Recomputed 是累计重新放置的网格数（规模观测）。
func (b *Batcher) Recomputed() int { return b.recomputed }

// Scanned 是累计扫描的网格数（规模观测）。
func (b *Batcher) Scanned() int { return b.scanned }
