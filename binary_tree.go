package lsmt

import (
  "fmt"
)

type TreeNode struct {
  Elem Element
  Left *TreeNode
  Right *TreeNode
  Size int
}

// NewTree accepts a sorted element slice and returns a binary tree representation.
func NewTree(elems []Element) *TreeNode {
  size := len(elems)
  if size == 0 {
    return nil
  }
  root := &TreeNode{
    Elem: elems[size/2],
    Left: NewTree(elems[0:size/2]),
    Size: size,
  }
  if rightIndex := size/2+1; rightIndex < size {
    root.Right = NewTree(elems[rightIndex:size])
  }
  return root
}

func Insert(tree **TreeNode, elem Element) {
  if *tree == nil {
    *tree = &TreeNode{Elem: elem}
  } else if elem.Key <= (*tree).Elem.Key {
    Insert(&((*tree).Left), elem)
  } else {
    Insert(&((*tree).Right), elem)
  }
  (*tree).Size++
}

func Find(tree *TreeNode, key string) (Element, error) {
  if tree == nil {
    // Not found.
    return Element{}, fmt.Errorf("key %s not found", key)
  } else if tree.Elem.Key == key {
    return tree.Elem, nil
  }
  if key <= tree.Elem.Key {
    return Find(tree.Left, key)
  } else {
    return Find(tree.Right, key)
  }
}

// Traverse returns all the elements in key order.
func Traverse(tree *TreeNode) []Element {
  var elems []Element
  if tree == nil {
    return elems
  }
  left := Traverse(tree.Left)
  right := Traverse(tree.Right)
  elems = append(elems, left...)
  elems = append(elems, tree.Elem)
  return append(elems, right...)
}

func JustSmallerOrEqual(tree *TreeNode, key string) (Element, error) {
  if tree == nil {
    return Element{}, fmt.Errorf("key %s is smaller than any key in the tree", key)
  }
  current := tree.Elem
  if current.Key <= key {
    right, err := JustSmallerOrEqual(tree.Right, key)
    if err == nil && current.Key < right.Key {
      current = right
    }
  } else {
    left, err := JustSmallerOrEqual(tree.Left, key)
    if err != nil {
      return Element{}, err
    }
    current = left
  }
  return current, nil
}

func JustLarger(tree *TreeNode, key string) (Element, error) {
  if tree == nil {
    return Element{}, fmt.Errorf("key %s is larger than any key in the tree", key)
  }
  current := tree.Elem
  if current.Key > key {
    left, err := JustLarger(tree.Left, key)
    if err == nil && current.Key > left.Key {
      current = left
    }
  } else {
    right, err := JustLarger(tree.Right, key)
    if err != nil {
      return Element{}, err
    }
    current = right
  }
  return current, nil
}