// Package batch 是网格合批：按材质分组、顶点上限拆批、加入顺序稳定、移除与改上限后重新紧凑化。
package batch

// Mesh 是一个待合批的网格。
type Mesh struct {
	ID       string
	Material string
	Verts    int
}

// chain 是同一材质下的网格序列与批次链。
type chain struct {
	order   []string       // 持有网格的加入顺序（含暂未入批的）
	verts   map[string]int // 网格顶点数
	batches [][]string     // 批次链，批内保持加入顺序
	loads   []int          // 每批的顶点总数
	where   map[string]int // 网格所在批次下标，-1 表示超上限暂未入批
}

func newChain() *chain {
	return &chain{
		verts: map[string]int{},
		where: map[string]int{},
	}
}

// place 把网格放进链尾：能进最后一批就进，否则开新批；floor 之前的批次不许动。
func (c *chain) place(id string, floor, maxVerts int) {
	v := c.verts[id]
	last := len(c.batches) - 1
	if last < floor || c.loads[last]+v > maxVerts {
		c.batches = append(c.batches, []string{id})
		c.loads = append(c.loads, v)
		last = len(c.batches) - 1
	} else {
		c.batches[last] = append(c.batches[last], id)
		c.loads[last] += v
	}
	c.where[id] = last
}

// Batcher 是合批器。
type Batcher struct {
	maxVerts   int
	chains     map[string]*chain
	matOrder   []string
	meshMat    map[string]string
	recomputed int
	scanned    int
}

// New 建合批器。
func New(maxVerts int) *Batcher {
	return &Batcher{
		maxVerts: maxVerts,
		chains:   map[string]*chain{},
		meshMat:  map[string]string{},
	}
}

// Add 加入一个网格。空材质或顶点数不大于 0 的网格直接忽略；
// 超过当前上限的网格暂不入批，等上限放宽后再按加入顺序放置。
func (b *Batcher) Add(m Mesh) {
	if m.Material == "" || m.Verts <= 0 {
		return
	}
	if _, dup := b.meshMat[m.ID]; dup {
		return
	}
	c := b.chains[m.Material]
	if c == nil {
		c = newChain()
		b.chains[m.Material] = c
		b.matOrder = append(b.matOrder, m.Material)
	}
	b.meshMat[m.ID] = m.Material
	c.order = append(c.order, m.ID)
	c.verts[m.ID] = m.Verts
	if m.Verts > b.maxVerts {
		c.where[m.ID] = -1
		return
	}
	c.place(m.ID, 0, b.maxVerts)
}

// Remove 移除一个网格，返回它原本是否存在。
func (b *Batcher) Remove(id string) bool {
	mat, ok := b.meshMat[id]
	if !ok {
		return false
	}
	delete(b.meshMat, id)
	c := b.chains[mat]
	pos := 0
	for pos < len(c.order) && c.order[pos] != id {
		pos++
	}
	c.order = append(c.order[:pos], c.order[pos+1:]...)
	delete(c.verts, id)
	batchIdx := c.where[id]
	delete(c.where, id)
	if batchIdx < 0 {
		return true
	}
	// 从被删批次的前一批起重排链尾，重新满足紧凑不变量。
	from := batchIdx - 1
	if from < 0 {
		from = 0
	}
	b.repackFrom(c, from)
	return true
}

// repackFrom 按加入顺序重排批次链从 from 开始的尾部，只碰这一条材质链。
func (b *Batcher) repackFrom(c *chain, from int) {
	pending := []string{}
	for _, ids := range c.batches[from:] {
		for _, id := range ids {
			if _, ok := c.where[id]; ok {
				pending = append(pending, id)
			}
		}
	}
	c.batches = c.batches[:from]
	c.loads = c.loads[:from]
	b.scanned += len(pending)
	for _, id := range pending {
		old := c.where[id]
		c.place(id, from, b.maxVerts)
		if c.where[id] != old {
			b.recomputed++
		}
	}
}

// SetLimit 改顶点上限，并按新上限重排所有批次链。
func (b *Batcher) SetLimit(limit int) {
	b.maxVerts = limit
	for _, mat := range b.matOrder {
		b.repackAll(b.chains[mat])
	}
}

// repackAll 按当前上限重排整条链：超上限的网格移出批次，其余按加入顺序入批。
func (b *Batcher) repackAll(c *chain) {
	c.batches = nil
	c.loads = nil
	b.scanned += len(c.order)
	for _, id := range c.order {
		old := c.where[id]
		if c.verts[id] > b.maxVerts {
			c.where[id] = -1
		} else {
			c.place(id, 0, b.maxVerts)
		}
		if c.where[id] != old {
			b.recomputed++
		}
	}
}

// Batches 返回每批的网格 id 列表，批次按材质首次出现的顺序排列。
func (b *Batcher) Batches() [][]string {
	out := [][]string{}
	for _, mat := range b.matOrder {
		for _, ids := range b.chains[mat].batches {
			out = append(out, append([]string{}, ids...))
		}
	}
	return out
}

// Recomputed 是累计重新放置的网格数（规模观测）。
func (b *Batcher) Recomputed() int { return b.recomputed }

// Scanned 是累计扫描的网格数（规模观测）。
func (b *Batcher) Scanned() int { return b.scanned }
