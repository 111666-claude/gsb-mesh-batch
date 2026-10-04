# mesh-batch

网格合批：按材质分组、顶点上限拆批、加入顺序稳定、移除与改上限后重新紧凑化，只用 Go 标准库。

```
go run ./cmd/mesh-batch --sample=material
go run ./cmd/mesh-batch --sample=limit
go run ./cmd/mesh-batch --sample=empty
go run ./cmd/mesh-batch --sample=remove
go run ./cmd/mesh-batch --sample=relimit
go run ./cmd/mesh-batch --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **材质分组**：只有材质相同的网格才能进同一批。
- **跳过**：材质为空、顶点数不大于 0、或顶点数大于 `maxVerts` 的网格直接忽略，不占批。
- **上限**：每批的顶点总数不超过 `maxVerts`。
- **顺序**：批内保持加入顺序；批次顺序按该材质首次出现的顺序。
- **保序**：同一材质里，加入更早的网格所在批次不晚于加入更晚的网格所在批次。
- **紧凑**：同一材质里，除最后一批外，任何一批再加上紧随其后的那个网格都会超上限。
- **移除与改上限**：`Remove` 之后该网格不再出现在任何批次里，`SetLimit` 之后按新上限重排；
  两种操作之后都必须重新满足上面全部不变量。
- **代价**：加入、移除与改上限只许碰受影响材质的批次链，`Scanned` 不随批次总数放大；
  `Recomputed` 只计被重新放置的网格。10 万个网格、每秒 10 万次操作。

## 输出契约（不改格式）

```
batches=..
meshes=..
ok=.. count=..
scanned=..
```

## 目录

```
batch.go              材质分组、上限、紧凑化与重排
batch_test.go         单元测试
cmd/mesh-batch/       命令行入口
```
