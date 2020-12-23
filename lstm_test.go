package lsmt

import (
  "fmt"
  "reflect"
  "sync"
  "testing"
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