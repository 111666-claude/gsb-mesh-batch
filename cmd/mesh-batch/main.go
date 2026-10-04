// Command mesh-batch 跑网格合批样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/batch"
)

func meshTotal(batcher *batch.Batcher) int {
	total := 0
	for _, item := range batcher.Batches() {
		total += len(item)
	}
	return total
}

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("mesh-batch", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "material", "material / limit / empty / remove / relimit / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *sample == "material" {
		batcher := batch.New(1000)
		batcher.Add(batch.Mesh{ID: "a", Material: "stone", Verts: 10})
		batcher.Add(batch.Mesh{ID: "b", Material: "wood", Verts: 10})
		fmt.Fprintf(stdout, "batches=%d\n", len(batcher.Batches()))
		return 0
	}
	if *sample == "limit" {
		batcher := batch.New(100)
		for _, id := range []string{"a", "b", "c"} {
			batcher.Add(batch.Mesh{ID: id, Material: "stone", Verts: 60})
		}
		fmt.Fprintf(stdout, "batches=%d\n", len(batcher.Batches()))
		return 0
	}
	if *sample == "empty" {
		batcher := batch.New(100)
		batcher.Add(batch.Mesh{ID: "a", Material: "", Verts: 0})
		batcher.Add(batch.Mesh{ID: "b", Material: "stone", Verts: 10})
		fmt.Fprintf(stdout, "meshes=%d\n", meshTotal(batcher))
		return 0
	}
	if *sample == "remove" {
		batcher := batch.New(100)
		for _, id := range []string{"a", "b", "c"} {
			batcher.Add(batch.Mesh{ID: id, Material: "stone", Verts: 10})
		}
		ok := batcher.Remove("b")
		fmt.Fprintf(stdout, "ok=%v count=%d\n", ok, meshTotal(batcher))
		return 0
	}
	if *sample == "relimit" {
		batcher := batch.New(1000)
		for _, id := range []string{"a", "b", "c"} {
			batcher.Add(batch.Mesh{ID: id, Material: "stone", Verts: 60})
		}
		batcher.SetLimit(100)
		fmt.Fprintf(stdout, "batches=%d\n", len(batcher.Batches()))
		return 0
	}
	if *sample == "scan" {
		batcher := batch.New(1000000)
		for index := 0; index < 3000; index++ {
			batcher.Add(batch.Mesh{
				ID:       fmt.Sprintf("m-%d", index),
				Material: fmt.Sprintf("mat-%d", index),
				Verts:    4,
			})
		}
		fmt.Fprintf(stdout, "scanned=%d\n", batcher.Scanned())
		return 0
	}
	fmt.Fprintln(stderr, "需要 --sample material|limit|empty|remove|relimit|scan")
	return 2
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
