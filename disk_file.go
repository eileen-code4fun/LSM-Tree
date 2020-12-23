package lsmt

import (
  "bytes"
  "encoding/gob"
  "fmt"
  "log"
  "io"
  "strconv"
)

const (
  maxFileLen = 1024
  indexSparseRatio = 3
)

type DiskFile struct {
  index *TreeNode
  data io.ReadSeeker
  size int
  buf bytes.Buffer
}

func (d DiskFile) Empty() bool {
  return d.size == 0
}

func NewDiskFile(elems []Element) DiskFile {
  d := DiskFile{size: len(elems)}
  var indexElems []Element
  var enc *gob.Encoder
  for i, e := range elems {
    if i % indexSparseRatio == 0 {
      // Create sparse index.
      idx := Element{Key: e.Key, Value: fmt.Sprintf("%d", d.buf.Len())}
      log.Printf("created sparse index element %v", idx)
      indexElems = append(indexElems, idx)
      enc = gob.NewEncoder(&d.buf)
    }
    enc.Encode(e)
  }
  d.index = NewTree(indexElems)
  return d
}

func (d DiskFile) Search(key string) (Element, error) {
  canErr := fmt.Errorf("key %s not found in disk file", key)
  if d.Empty() {
    return Element{}, canErr
  }
  var si, ei int
  start, err := JustSmallerOrEqual(d.index, key)
  if err != nil {
    // Key smaller than all.
    return Element{}, canErr
  }
  si, _ = strconv.Atoi(start.Value)
  end, err := JustLarger(d.index, key)
  if err != nil {
    // Key larger than all or equal to the last one.
    ei = d.buf.Len()
  } else {
    ei, _ = strconv.Atoi(end.Value)
  }
  log.Printf("searching in range [%d,%d)]", si, ei)
  buf := bytes.NewBuffer(d.buf.Bytes()[si:ei])
  dec := gob.NewDecoder(buf)
  for {
    var e Element
    if err := dec.Decode(&e); err != nil {
      log.Printf("got err: %v", err)
      break
    }
    if e.Key == key {
      return e, nil
    }
  }
  return Element{}, canErr
}

func (d DiskFile) AllElements() []Element {
  indexElems := Traverse(d.index)
  var elems []Element
  var dec *gob.Decoder
  for i, idx := range indexElems {
    start, _ := strconv.Atoi(idx.Value)
    end := d.buf.Len()
    if i < len(indexElems)-1 {
      end, _ = strconv.Atoi(indexElems[i+1].Value)
    }
    dec = gob.NewDecoder(bytes.NewBuffer(d.buf.Bytes()[start:end]))
    var e Element
    for dec.Decode(&e)==nil {
      elems = append(elems, e)
    }
  }
  return elems
}