package tree

import (
	"reflect"
	"testing"
)

type Menu struct {
	ID       int
	ParentID int
	SubMenus []*Menu
}

// TestMakeTree MakeTree单测，从列表构建树
func TestMakeTree(t *testing.T) {
	menus := []*Menu{
		{ID: 1, ParentID: 0},
		{ID: 2, ParentID: 1},
		{ID: 3, ParentID: 1},
		{ID: 4, ParentID: 2},
	}

	isRoot := func(menu *Menu) bool {
		return menu.ParentID == 0
	}

	isParent := func(parent, child *Menu) bool {
		return parent.ID == child.ParentID
	}

	setSubMenus := func(parent *Menu, child []*Menu) {
		parent.SubMenus = child
	}

	expected := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
			{ID: 3, ParentID: 1},
		}},
	}

	tree := MakeTree(menus, isRoot, isParent, setSubMenus)
	if !reflect.DeepEqual(tree, expected) {
		t.Errorf("expected %v, got %v", expected, tree)
	}
}

// TestFlat Flat单测，将树打平为列表
func TestFlat(t *testing.T) {

	tree := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2, SubMenus: []*Menu{}},
			}},
			{ID: 3, ParentID: 1, SubMenus: []*Menu{}},
		}},
	}

	getSubMenus := func(menu *Menu) []*Menu {
		return menu.SubMenus
	}

	clearSubMenus := func(menu *Menu) {
		menu.SubMenus = nil
	}

	expected := []*Menu{
		{ID: 4, ParentID: 2, SubMenus: nil},
		{ID: 2, ParentID: 1, SubMenus: nil},
		{ID: 3, ParentID: 1, SubMenus: nil},
		{ID: 1, ParentID: 0, SubMenus: nil},
	}

	flat := Flat(tree, getSubMenus, clearSubMenus)
	convert := func(menus []*Menu) []Menu {
		result := make([]Menu, len(menus))
		for i, menu := range menus {
			result[i] = *menu
		}
		return result
	}

	if !reflect.DeepEqual(convert(flat), convert(expected)) {
		t.Errorf("expected %v, got %v", convert(expected), convert(flat))
	}
}

// TestForPreOrder ForPreOrder单测，前序遍历
func TestForPreOrder(t *testing.T) {

	tree := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
			{ID: 3, ParentID: 1},
		}},
	}

	getSubMenus := func(menu *Menu) []*Menu {
		return menu.SubMenus
	}

	var ids []int
	consumer := func(menu *Menu) {
		ids = append(ids, menu.ID)
	}

	expected := []int{1, 2, 4, 3}

	ForPreOrder(tree, consumer, getSubMenus)
	if !reflect.DeepEqual(ids, expected) {
		t.Errorf("expected %v, got %v", expected, ids)
	}
}

// TestForLevelOrder ForLevelOrder单测，层序遍历
func TestForLevelOrder(t *testing.T) {

	tree := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
			{ID: 3, ParentID: 1},
		}},
	}

	getSubMenus := func(menu *Menu) []*Menu {
		return menu.SubMenus
	}

	var ids []int
	consumer := func(menu *Menu) {
		ids = append(ids, menu.ID)
	}

	expected := []int{1, 2, 3, 4}

	ForLevelOrder(tree, consumer, getSubMenus)
	if !reflect.DeepEqual(ids, expected) {
		t.Errorf("expected %v, got %v", expected, ids)
	}
}

// TestForPostOrder ForPostOrder单测，后序遍历
func TestForPostOrder(t *testing.T) {

	tree := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
			{ID: 3, ParentID: 1},
		}},
	}

	getSubMenus := func(menu *Menu) []*Menu {
		return menu.SubMenus
	}

	var ids []int
	consumer := func(menu *Menu) {
		ids = append(ids, menu.ID)
	}

	expected := []int{4, 2, 3, 1}

	ForPostOrder(tree, consumer, getSubMenus)
	if !reflect.DeepEqual(ids, expected) {
		t.Errorf("expected %v, got %v", expected, ids)
	}
}

// TestSort Sort单测，对树的所有子节点进行排序
func TestSort(t *testing.T) {
	tree := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 3, ParentID: 1},
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
		}},
	}

	comparator := func(a, b *Menu) bool {
		return a.ID < b.ID
	}

	getChildren := func(menu *Menu) []*Menu {
		return menu.SubMenus
	}

	expected := []*Menu{
		{ID: 1, ParentID: 0, SubMenus: []*Menu{
			{ID: 2, ParentID: 1, SubMenus: []*Menu{
				{ID: 4, ParentID: 2},
			}},
			{ID: 3, ParentID: 1},
		}},
	}

	sorted := Sort(tree, comparator, getChildren)
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("expected %v, got %v", expected, sorted)
	}
}
