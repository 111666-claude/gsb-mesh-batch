// Command mesh-batch 跑网格合批样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/batch"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("mesh-batch", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "material", "material / limit / empty / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "material":
		batcher := batch.New(1000)
		batcher.Add(batch.Mesh{ID: "a", Material: "stone", Verts: 10})
		batcher.Add(batch.Mesh{ID: "b", Material: "wood", Verts: 10})
		fmt.Fprintf(stdout, "batches=%d\n", len(batcher.Batches()))
	case "limit":
		batcher := batch.New(100)
		for _, id := range []string{"a", "b", "c"} {
			batcher.Add(batch.Mesh{ID: id, Material: "stone", Verts: 60})
		}
		fmt.Fprintf(stdout, "batches=%d\n", len(batcher.Batches()))
	case "empty":
		batcher := batch.New(100)
		batcher.Add(batch.Mesh{ID: "a", Material: "", Verts: 0})
		batcher.Add(batch.Mesh{ID: "b", Material: "stone", Verts: 10})
		batches := batcher.Batches()
		total := 0
		for _, item := range batches {
			total += len(item)
		}
		fmt.Fprintf(stdout, "meshes=%d\n", total)
	case "scan":
		batcher := batch.New(1000000)
		for index := 0; index < 3000; index++ {
			batcher.Add(batch.Mesh{ID: fmt.Sprintf("m-%d", index), Material: "stone", Verts: 4})
		}
		fmt.Fprintf(stdout, "scanned=%d\n", batcher.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample material|limit|empty|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
