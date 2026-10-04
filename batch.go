// Package batch 是网格合批：按材质分组、顶点上限拆批、加入顺序稳定、移除与改上限后重新紧凑化。
//
// 每个材质维护一条独立的批次链；加入时只看该材质最后一批的顶点和，
// 移除只重排受影响材质的链，改上限时对每条链按新上限重新贪心放置。
package batch

// Mesh 是一个待合批的网格。
type Mesh struct {
	ID       string
	Material string
	Verts    int
}

// matChain 是单个材质的网格链与批次表；批次里存的是 meshes 的下标。
type matChain struct {
	meshes  []Mesh
	batches [][]int
	sums    []int
}

// Batcher 是合批器。
type Batcher struct {
	maxVerts   int
	chains     map[string]*matChain
	order      []string
	materialOf map[string]string
	recomputed int
	scanned    int
}

// New 建合批器。
func New(maxVerts int) *Batcher {
	return &Batcher{
		maxVerts:   maxVerts,
		chains:     map[string]*matChain{},
		materialOf: map[string]string{},
	}
}

// Add 加入一个网格：材质为空、顶点数不大于 0 或超过 maxVerts 的直接忽略。
func (b *Batcher) Add(m Mesh) {
	b.scanned++
	if m.Material == "" || m.Verts <= 0 || m.Verts > b.maxVerts {
		return
	}
	chain, ok := b.chains[m.Material]
	if !ok {
		chain = &matChain{}
		b.chains[m.Material] = chain
		b.order = append(b.order, m.Material)
	}
	idx := len(chain.meshes)
	chain.meshes = append(chain.meshes, m)
	if last := len(chain.batches) - 1; last >= 0 && chain.sums[last]+m.Verts <= b.maxVerts {
		chain.batches[last] = append(chain.batches[last], idx)
		chain.sums[last] += m.Verts
	} else {
		chain.batches = append(chain.batches, []int{idx})
		chain.sums = append(chain.sums, m.Verts)
	}
	b.materialOf[m.ID] = m.Material
}

// Remove 移除一个网格，返回它原本是否存在；只重排该材质的批次链。
func (b *Batcher) Remove(id string) bool {
	material, ok := b.materialOf[id]
	if !ok {
		return false
	}
	chain := b.chains[material]
	pos := -1
	for i, m := range chain.meshes {
		b.scanned++
		if m.ID == id {
			pos = i
			break
		}
	}
	if pos < 0 {
		return false
	}
	old := b.snapshot(chain)
	chain.meshes = append(chain.meshes[:pos], chain.meshes[pos+1:]...)
	delete(b.materialOf, id)
	b.repack(chain, old)
	return true
}

// SetLimit 改顶点上限，并按新上限重排全部材质链。
func (b *Batcher) SetLimit(limit int) {
	b.maxVerts = limit
	for _, material := range b.order {
		chain := b.chains[material]
		b.repack(chain, b.snapshot(chain))
	}
}

// snapshot 记录链上每个网格当前所在的批次号。
func (b *Batcher) snapshot(chain *matChain) map[string]int {
	old := make(map[string]int, len(chain.meshes))
	for bi, batch := range chain.batches {
		for _, idx := range batch {
			old[chain.meshes[idx].ID] = bi
		}
	}
	return old
}

// repack 按当前 maxVerts 贪心重排一条材质链；
// scanned 只记本次碰到的链上网格，recomputed 只记批次位置发生变化的网格。
func (b *Batcher) repack(chain *matChain, old map[string]int) {
	chain.batches = nil
	chain.sums = nil
	var current []int
	sum := 0
	for idx, m := range chain.meshes {
		b.scanned++
		newBatch := -1
		if m.Verts <= b.maxVerts {
			if sum > 0 && sum+m.Verts > b.maxVerts {
				chain.batches = append(chain.batches, current)
				chain.sums = append(chain.sums, sum)
				current = nil
				sum = 0
			}
			newBatch = len(chain.batches)
			current = append(current, idx)
			sum += m.Verts
		}
		prev, ok := old[m.ID]
		if !ok {
			prev = -1
		}
		if prev != newBatch {
			b.recomputed++
		}
	}
	if len(current) > 0 {
		chain.batches = append(chain.batches, current)
		chain.sums = append(chain.sums, sum)
	}
}

// Batches 返回每批的网格 id 列表（深拷贝），批次按材质首次出现顺序排列。
func (b *Batcher) Batches() [][]string {
	out := [][]string{}
	for _, material := range b.order {
		for _, batch := range b.chains[material].batches {
			ids := make([]string, len(batch))
			for i, idx := range batch {
				ids[i] = b.chains[material].meshes[idx].ID
			}
			out = append(out, ids)
		}
	}
	return out
}

// Recomputed 是累计重新放置的网格数（规模观测）。
func (b *Batcher) Recomputed() int { return b.recomputed }

// Scanned 是累计扫描的网格数（规模观测）。
func (b *Batcher) Scanned() int { return b.scanned }
