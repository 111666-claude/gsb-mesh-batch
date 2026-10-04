# mesh-batch

网格合批：按材质分组、顶点上限拆批与顺序稳定，只用 Go 标准库。

```
go run ./cmd/mesh-batch --sample=material
go run ./cmd/mesh-batch --sample=limit
go run ./cmd/mesh-batch --sample=empty
go run ./cmd/mesh-batch --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **材质分组**：只有材质相同的网格才能进同一批。
- **顶点上限**：一批的顶点总数不超过 `maxVerts`，放不下就开新批。
- **顺序**：批内保持加入顺序，批次顺序按该材质首次出现的顺序。
- **跳过**：材质为空或顶点数不大于 0 的网格直接忽略，不占批。
- **不变量**：每批顶点数不超过上限；同一串加入序列重复执行结果相同。
- **代价**：加入网格不许重扫全部已有批次，`Scanned` 不随批次规模平方增长。
- **规模**：10 万个网格、每秒 10 万次加入，单次摊还 O(1)。

## 输出契约（不改格式）

```
batches=..
meshes=..
scanned=..
```

## 目录

```
batch.go              材质分组、上限与跳过
batch_test.go         单元测试
cmd/mesh-batch/       命令行入口
```
