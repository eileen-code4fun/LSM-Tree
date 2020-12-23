package lsmt

import (
  "fmt"
  "reflect"
  "sync"
  "testing"
  "time"
)

func TestInMemoryOnly(t *testing.T) {
  var wg sync.WaitGroup
  var expected []Element
  total := 10
  tree := NewLSMTree(total+1 /* flush threshold larger than total */)
  for i := 0; i < total; i ++ {
    e := Element{Key: fmt.Sprintf("%d", i), Value: fmt.Sprintf("%d", i)}
    expected = append(expected, e)
    wg.Add(1)
    go func(){
      tree.Put(e.Key, e.Value)
      wg.Done()
    }()
  }
  wg.Wait()
  if tree.tree.Size != total {
    t.Errorf("got tree size %d; want %d", tree.tree.Size, total)
  }
  for i := 0; i < total; i ++ {
    wg.Add(1)
    e := fmt.Sprintf("%d", i)
    go func() {
      v, err := tree.Get(e)
      if err != nil {
        t.Errorf("key %s not found", e)
      }
      if v != e {
        t.Errorf("got %s for key %s; want %s", v, e, e)
      }
      wg.Done()
    }()
  }
  wg.Wait()
  got := Traverse(tree.tree)
  if !reflect.DeepEqual(expected, got) {
    t.Errorf("got result %v; want %v", got, expected)
  }
}

func TestFlushedToDisk(t *testing.T) {
  t.Parallel()
  tree := NewLSMTree(2)
  tree.Put("1", "One")
  tree.Put("2", "Two")
  // Wait for flush.
  time.Sleep(1 * time.Second)
  if tree.tree != nil {
    t.Errorf("got tree size %d; want empty", tree.tree.Size)
  }
  if len(tree.diskFiles) != 1 {
    t.Errorf("got disk file size %d; want 1", len(tree.diskFiles))
  }
  if _, err := tree.Get("1"); err != nil {
      t.Error("key 1 not found")
    }
  if _, err := tree.Get("2"); err != nil {
    t.Error("key 2 not found")
  }
  tree.Put("3", "Three")
  if _, err := tree.Get("3"); err != nil {
    t.Error("key 3 not found")
  }
  tree.Put("4", "Four")
  // Wait for flush and compaction.
  time.Sleep(3 * time.Second)
  if len(tree.diskFiles) != 1 {
    t.Errorf("got disk file size %d; want 1", len(tree.diskFiles))
  }
  if len(tree.diskFiles) == 1 {
    got := tree.diskFiles[0].AllElements()
    want := []Element{{Key: "1", Value: "One"}, {Key: "2", Value: "Two"}, {Key: "3", Value: "Three"}, {Key: "4", Value: "Four"}}
    if !reflect.DeepEqual(want, got) {
      t.Errorf("got result %v; want %v", got, want)
    }
  }
}

func TestCompactionCollapse(t *testing.T) {
  t.Parallel()
  tree := NewLSMTree(1)
  tree.Put("1", "One")
  time.Sleep(time.Second)
  tree.Put("1", "ONE")
  // Wait for flush and compaction.
  time.Sleep(3 * time.Second)
  if len(tree.diskFiles) != 1 {
    t.Errorf("got disk file size %d; want 1", len(tree.diskFiles))
  }
  if len(tree.diskFiles) == 1 {
    got := tree.diskFiles[0].AllElements()
    want := []Element{{Key: "1", Value: "ONE"}}
    if !reflect.DeepEqual(want, got) {
      t.Errorf("got result %v; want %v", got, want)
    }
  }
}