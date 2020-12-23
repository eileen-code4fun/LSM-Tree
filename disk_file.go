package lsmt

import (
  "fmt"
  "io"
)

type DiskFile struct {
  index *TreeNode
  data io.ReadSeeker
  size int
  buf []byte
}

func (d DiskFile) Empty() bool {
  return d.size == 0
}

func NewDiskFile(elems []Element) DiskFile {
  // TODO
  return DiskFile{size: len(elems)}
}

func (d DiskFile) Search(key string) (string, error) {
  // TODO
  return "", fmt.Errorf("key %s not found in disk file", key)
}

func (d DiskFile) AllElements() []Element {
  // TODO
  return nil
}