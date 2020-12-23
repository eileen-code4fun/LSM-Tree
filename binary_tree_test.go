package lsmt

import (
  "reflect"
  "testing"
)

func TestBinaryTree(t *testing.T) {
  elems := []Element{
    {Key: "1", Value: "One"},
    {Key: "2", Value: "Two"},
    {Key: "3", Value: "Three"},
    {Key: "4", Value: "Four"},
    {Key: "5", Value: "Five"},
    {Key: "6", Value: "Six"},
    {Key: "7", Value: "Seven"},
  }
  tree := NewTree(elems)
  if tree.Size != len(elems) {
    t.Errorf("got tree size %d; want %d", tree.Size, len(elems))
  }
  e, err := Find(tree, "5")
  if err != nil {
    t.Errorf("find 5 got error: %v", err)
  }
  if e.Key != "5" || e.Value != "Five" {
    t.Errorf("got key %s value %s; want 5 Five", e.Key, e.Value)
  }
  e, err = Find(tree, "2")
  if err != nil {
    t.Errorf("find 2 got error: %v", err)
  }
  if e.Key != "2" || e.Value != "Two" {
    t.Errorf("got key %s value %s; want 2 Two", e.Key, e.Value)
  }
  if e, err := Find(tree, "1.5"); err == nil {
    t.Errorf("got elem %v; want not found", e)
  }
  newElem := Element{Key: "6.5", Value: "Six Point Five"}
  Insert(&tree, newElem)
  got := Traverse(tree)
  var expected []Element
  expected = append(expected, elems[0:6]...)
  expected = append(expected, newElem)
  expected = append(expected, elems[6:]...)
  if tree.Size != len(expected) {
    t.Errorf("got tree size %d; want %d", tree.Size, len(expected))
  }
  if !reflect.DeepEqual(expected, got) {
    t.Errorf("traverse got %v; want %v", got, expected)
  }
}

func TestTreeKeyComparison(t *testing.T) {
  elems := []Element{
    {Key: "1", Value: "One"},
    {Key: "2", Value: "Two"},
    {Key: "3", Value: "Three"},
    {Key: "4", Value: "Four"},
    {Key: "5", Value: "Five"},
    {Key: "6", Value: "Six"},
    {Key: "7", Value: "Seven"},
  }
  tree := NewTree(elems)
  if e, err := JustSmallerOrEqual(tree, "1"); err != nil || e.Key != "1" {
    t.Errorf("got %v, %v; want key 1, nil", e, err)
  }
  if e, err := JustSmallerOrEqual(tree, "0"); err == nil {
    t.Errorf("got %v, %v; want not found", e, err)
  }
  if e, err := JustSmallerOrEqual(tree, "2.5"); err != nil || e.Key != "2" {
    t.Errorf("got %v, %v; want key 2, nil", e, err)
  }
  if e, err := JustSmallerOrEqual(tree, "8"); err != nil || e.Key != "7" {
    t.Errorf("got %v, %v; want key 7, nil", e, err)
  }
  if e, err := JustLarger(tree, "6"); err != nil || e.Key != "7" {
    t.Errorf("got %v, %v; want key 7, nil", e, err)
  }
  if e, err := JustLarger(tree, "7"); err == nil {
    t.Errorf("got %v, %v; want not found", e, err)
  }
  if e, err := JustLarger(tree, "0"); err != nil || e.Key != "1" {
    t.Errorf("got %v, %v; want key 1, nil", e, err)
  }
}