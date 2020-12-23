package lsmt

import (
  "fmt"
  "sync"
  "time"
)

type Element struct {
  Key, Value string
}

type LSMTree struct {
  // Read write lock to control access to the in-memory tree.
  rwm sync.RWMutex
  tree *TreeNode
  treeInFlush *TreeNode
  flushThreshold int
  // Read write lock to control access to the disk files.
  drwm sync.RWMutex
  diskFiles []DiskFile
}

func NewLSMTree(flushThreshold int) *LSMTree {
  t := &LSMTree{flushThreshold: flushThreshold}
  go t.compactService()
  return t
}

func (t *LSMTree) Put(key, value string) {
  t.rwm.Lock()
  defer t.rwm.Unlock()
  Insert(&(t.tree), Element{Key: key, Value: value})
  if t.tree.Size >= t.flushThreshold && t.treeInFlush == nil {
    // Trigger flush.
    t.treeInFlush = t.tree
    t.tree = nil
    go t.flush()
  }
}

func (t *LSMTree) Get(key string) (string, error) {
  t.rwm.RLock()
  if e, err := Find(t.tree, key); err == nil {
    t.rwm.RUnlock()
    return e.Value, nil
  }
  if e, err := Find(t.treeInFlush, key); err == nil {
    t.rwm.RUnlock()
    return e.Value, nil
  }
  t.rwm.RUnlock()
  // The key is not in memory. Search in disk files.
  t.drwm.RLock()
  defer t.drwm.RUnlock()
  for _, d := range t.diskFiles {
    e, err := d.Search(key)
    if err == nil {
      // Found in disk
      return e.Value, nil
    }
  }
  return "", fmt.Errorf("key %s not found", key)
}

func (t *LSMTree) flush() {
  // Create a new disk file.
  d := []DiskFile{NewDiskFile(Traverse(t.treeInFlush))}
  // Put the disk file in the list.
  t.drwm.Lock()
  t.diskFiles = append(d, t.diskFiles...)
  t.drwm.Unlock()
  // Remove the tree in flush.
  t.rwm.Lock()
  t.treeInFlush = nil
  t.rwm.Unlock()
}

func (t *LSMTree) compactService() {
  for {
    time.Sleep(time.Second)
    var d1, d2 DiskFile
    t.drwm.RLock()
    if len(t.diskFiles) >= 2 {
      d1 = t.diskFiles[len(t.diskFiles)-1]
      d2 = t.diskFiles[len(t.diskFiles)-2]
    }
    t.drwm.RUnlock()
    if d1.Empty() || d2.Empty() {
      continue
    }
    // Create a new compacted disk file.
    d := compact(d1, d2)
    // Replace the two old files.
    t.drwm.Lock()
    t.diskFiles = t.diskFiles[0:len(t.diskFiles)-2]
    t.diskFiles = append(t.diskFiles, d)
    t.drwm.Unlock()
  }
}

func compact(d1, d2 DiskFile) DiskFile {
  elems1 := d1.AllElements()
  elems2 := d2.AllElements()
  size := min(len(elems1), len(elems2))
  var newElems []Element
  var i1, i2 int
  for i1 < size && i2 < size {
    e1 := elems1[i1]
    e2 := elems2[i2]
    if e1.Key < e2.Key {
      newElems = append(newElems, e1)
      i1++
    } else {
      newElems = append(newElems, e2)
      i2++
    }
  }
  newElems = append(newElems, elems1[i1:len(elems1)]...)
  newElems = append(newElems, elems2[i2:len(elems2)]...)
  return NewDiskFile(newElems)
}

func min(i, j int) int {
  if i < j {
    return i
  }
  return j
}