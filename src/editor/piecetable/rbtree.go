package piecetable

//define enum node color
type NodeColor int

const (
	Red NodeColor = iota
	Black
)

type TreeNode struct {
	parent *TreeNode
	left   *TreeNode
	right  *TreeNode
	color  NodeColor

	piece     *Piece
	size_left uint
	lf_left   uint
}

func Constructor(piece *Piece, color NodeColor) (tree *TreeNode) {

	return &TreeNode{
		parent:    tree,
		left:      tree,
		right:     tree,
		color:     color,
		piece:     piece,
		size_left: 0,
		lf_left:   0,
	}
}

func (tree *TreeNode) Next() TreeNode {
	if tree.right != SENTINEL {
		return Leftest(tree.right)
	}
	node := tree

	for {
		if node.parent != SENTINEL {
			if node.parent.left == node {
				break
			}

			node = node.parent

		} else {
			break
		}
	}

	if node.parent == SENTINEL {
		return *SENTINEL
	} else {
		return *node.parent
	}

}

func (tree *TreeNode) Prev() TreeNode {
	if tree.left != SENTINEL {
		return Righttest(tree.left)
	}
	node := tree

	for {
		if node.parent != SENTINEL {
			if node.parent.right == node {
				break
			}

			node = node.parent

		} else {
			break
		}
	}

	if node.parent == SENTINEL {
		return *SENTINEL
	} else {
		return *node.parent
	}
}

func (tree *TreeNode) Detach() {
	tree.parent = nil
	tree.left = nil
	tree.right = nil
}

//instantiation one const tree node
var SENTINEL_TEMP *TreeNode = &TreeNode{
	parent:    nil,
	left:      nil,
	right:     nil,
	color:     NodeColor(Black),
	piece:     nil,
	size_left: 0,
	lf_left:   0,
}

func initSENTINEL(tree *TreeNode) *TreeNode {
	tree.parent = SENTINEL_TEMP
	tree.left = SENTINEL_TEMP
	tree.right = SENTINEL_TEMP
	tree.color = NodeColor(Black)

	return tree
}

// var SENTINEL = Constructor(nil, NodeColor(Black))

var SENTINEL *TreeNode = initSENTINEL(SENTINEL_TEMP)

//test node->left's color
func Leftest(node *TreeNode) TreeNode {
	for {
		if node.left != SENTINEL {
			node = node.left
		} else {
			break
		}
	}
	return *node
}

//test node->right's color
func Righttest(node *TreeNode) TreeNode {
	for {
		if node.right != SENTINEL {
			node = node.right
		} else {
			break
		}
	}
	return *node
}

func CalculateSize(node *TreeNode) uint {
	if node == SENTINEL {
		return 0
	}
	return node.size_left + node.piece.length + CalculateSize(node.right)
}

func CalculateLF(node *TreeNode) uint {
	if node == SENTINEL {
		return 0
	}
	return node.lf_left + node.piece.lineFeedCnt + CalculateLF(node.right)
}

func resetSentinel() {
	SENTINEL.parent = SENTINEL
}

func UpdateTreeMetadata(tree *PieceTreeBase, x *TreeNode, delta uint, lineFeedCntDelta uint) {
	for {
		if x != tree.root && x != SENTINEL {
			if x.parent.left == x {
				x.parent.size_left += delta
				x.parent.lf_left += lineFeedCntDelta
			}
		} else {
			break
		}
	}
}

func RecomputeTreeMetadata(tree *PieceTreeBase, x *TreeNode) {
	var delta uint = 0
	var lf_delta uint = 0

	if x == tree.root {
		return
	}

	//go upwards til the node whose left subtree is changed
	for {
		if x != tree.root && x == x.parent.right {
			x = x.parent
		} else {
			break
		}
	}

	//means we add a node to the end
	if x == tree.root {
		return
	}

	//x is the ndoe whose right subtree is changed
	x = x.parent

	delta = CalculateSize(x.left) - x.size_left
	lf_delta = CalculateLF(x.left) - x.lf_left

	x.size_left += delta
	x.lf_left += lf_delta

	//go upwards till root
	for {
		if x != tree.root && (delta != 0 || lf_delta != 0) {
			if x.parent.left == x {
				x.parent.size_left += delta
				x.parent.lf_left += lf_delta
			}
			x = x.parent
		} else {
			break
		}
	}
}

func LeftRotate(tree *PieceTreeBase, x *TreeNode) {
	var y = x.right
	var xPieceLength uint = 0
	var xPieceLineFeedCnt uint = 0

	if x.piece != nil {
		xPieceLength = x.piece.length
		xPieceLineFeedCnt = x.piece.lineFeedCnt
	}

	y.size_left += x.size_left + xPieceLength
	y.lf_left += x.lf_left + xPieceLineFeedCnt
	x.right = y.left

	if y.left != SENTINEL {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == SENTINEL {
		tree.root = y
	} else if x.parent.left == x {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y

}

func RightRotate(tree *PieceTreeBase, y *TreeNode) {
	var xPieceLength uint = 0
	var xPieceLineFeedCnt uint = 0
	var x = y.left

	y.left = x.right

	if x.right != SENTINEL {
		x.right.parent = y
	}

	x.parent = y.parent

	if x.piece != nil {
		xPieceLength = x.piece.length
		xPieceLineFeedCnt = x.piece.lineFeedCnt
	}

	y.size_left -= x.size_left + xPieceLength
	y.lf_left -= x.lf_left + xPieceLineFeedCnt

	if y.parent != SENTINEL {
		tree.root = x
	} else if y == y.parent.right {
		y.parent.right = x
	} else {
		y.parent.left = x
	}

	x.right = y
	y.parent = x

}

func RBDelete(tree *PieceTreeBase, z *TreeNode) {
	var x *TreeNode
	var y *TreeNode

	if z.left == SENTINEL {
		y = z
		x = y.right
	} else if z.right == SENTINEL {
		y = z
		x = y.left
	} else {
		*y = Leftest(z.right)
		x = y.right
	}

	if y == tree.root {
		tree.root = x
		x.color = NodeColor(Black)
		z.Detach()
		resetSentinel()
		tree.root.parent = SENTINEL

		return
	}

	var yWasRed = (y.color == NodeColor(Red))

	if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y == z {
		x.parent = y.parent
		RecomputeTreeMetadata(tree, x)
	} else {
		if y.parent == z {
			x.parent = y
		} else {
			x.parent = y.parent
		}

		RecomputeTreeMetadata(tree, x)

		y.left = z.left
		y.right = z.right
		y.parent = z.parent
		y.color = z.color

		if z == tree.root {
			tree.root = y
		} else {
			if z == z.parent.left {
				z.parent.left = y
			} else {
				z.parent.right = y
			}
		}

		if y.left != SENTINEL {
			y.left.parent = y
		}
		if y.right != SENTINEL {
			y.right.parent = y
		}

		y.size_left = z.size_left
		y.lf_left = z.lf_left
		RecomputeTreeMetadata(tree, y)
	}

	z.Detach()

	if x.parent.left == x {
		newSizeLeft := CalculateSize(x)
		newLFLeft := CalculateLF(x)

		if newLFLeft != x.parent.size_left || newLFLeft != x.parent.lf_left {
			delta := newSizeLeft - x.parent.size_left
			lf_delta := newLFLeft - x.parent.lf_left

			x.parent.size_left = newSizeLeft
			x.parent.lf_left = newLFLeft

			UpdateTreeMetadata(tree, x.parent, delta, lf_delta)
		}
	}

	RecomputeTreeMetadata(tree, x.parent)

	if yWasRed {
		resetSentinel()
		return
	}

	var w *TreeNode
	for {
		if x != tree.root && x.color == NodeColor(Black) {
			if x == x.parent.left {
				w = x.parent.right
				if w.color == NodeColor(Red) {
					w.color = NodeColor(Black)
					x.parent.color = NodeColor(Red)
					LeftRotate(tree, x.parent)
					w = x.parent.right
				}

				if w.left.color == NodeColor(Black) && w.right.color == NodeColor(Black) {
					w.color = NodeColor(Red)
					x = x.parent
				} else {
					if w.right.color == NodeColor(Black) {
						w.left.color = NodeColor(Black)
						w.color = NodeColor(Red)
						RightRotate(tree, w)
						w = x.parent.right
					}

					w.color = x.parent.color
					x.parent.color = NodeColor(Black)
					x.right.color = NodeColor(Black)
					LeftRotate(tree, x.parent)
					x = tree.root
				}
			} else {
				w = x.parent.left

				if w.color == NodeColor(Red) {
					w.color = NodeColor(Black)
					x.parent.color = NodeColor(Red)
					RightRotate(tree, x.parent)
					w = x.parent.left
				}

				if w.left.color == NodeColor(Black) && w.right.color == NodeColor(Black) {
					w.color = NodeColor(Red)
					x = x.parent
				} else {
					if w.left.color == NodeColor(Black) {
						w.right.color = NodeColor(Black)
						w.color = NodeColor(Red)
						LeftRotate(tree, w)
						w = x.parent.left
					}

					w.color = x.parent.color
					x.parent.color = NodeColor(Black)
					w.left.color = NodeColor(Black)
					RightRotate(tree, x.parent)
					x = tree.root
				}
			}

		} else {
			break
		}
	}
	x.color = NodeColor(Black)
	resetSentinel()

}

func FixInsert(tree *PieceTreeBase, x *TreeNode) {
	RecomputeTreeMetadata(tree, x)
	for {
		if x != tree.root && x.parent.color == NodeColor(Red) {
			if x.parent == x.parent.parent.left {
				y := x.parent.parent.right

				if y.color == NodeColor(Red) {
					x.parent.color = NodeColor(Black)
					y.color = NodeColor(Black)
					x.parent.parent.color = NodeColor(Red)
					x = x.parent.parent
				} else {
					if x == x.parent.right {
						x = x.parent
						LeftRotate(tree, x)
					}

					x.parent.color = NodeColor(Black)
					x.parent.parent.color = NodeColor(Red)
					RightRotate(tree, x.parent.parent)
				}
			} else {
				y := x.parent.parent.left

				if y.color == NodeColor(Red) {
					x.parent.color = NodeColor(Black)
					y.color = NodeColor(Black)
					x.parent.parent.color = NodeColor(Red)
					x = x.parent.parent
				} else {
					if x == x.parent.left {
						x = x.parent
						RightRotate(tree, x)
					}
					x.parent.color = NodeColor(Black)
					x.parent.parent.color = NodeColor(Red)
					LeftRotate(tree, x.parent.parent)
				}
			}

		} else {
			break
		}
	}
	tree.root.color = NodeColor(Black)
}
