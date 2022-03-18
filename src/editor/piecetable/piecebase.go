package piecetable

const AverageBufferSize = 65535

type Piece struct {
	length      uint
	lineFeedCnt uint
}

type PieceTreeBase struct {
	root *TreeNode
}

type LineStarts struct {
	lineStarts   uint
	cr           uint
	lf           uint
	crlf         uint
	isBasicASCII bool
}

func CreateUintArray[T ~uint](arr []T) interface{} {

	if arr[len(arr)-1] < AverageBufferSize {
		var r = make([]uint16, 0)

		return r

	} else {
		var r = make([]uint32, 0)
		return r
	}

}

//constructor of line starts
func ConstructorLineStarts(lineStarts uint, cr, lf, crlf uint, isBasicASCII bool) *LineStarts {
	return &LineStarts{
		lineStarts:   lineStarts,
		cr:           cr,
		lf:           lf,
		crlf:         crlf,
		isBasicASCII: isBasicASCII,
	}
}

func CreateLineStartsFast(str string, readonly bool) interface{} {
	readonly = true

	r := []uint{}
	var rLength uint = 1

	for i := 0; i < len(str); i++ {
		chr := rune(str[i])

		if chr == 13 {
			if i+1 < len(str) && rune(str[i+1]) == 10 {
				r[rLength] = uint(i) + 2
				i++
				rLength++
			} else {
				r[rLength] = uint(i) + 1
				rLength++
			}
		} else if chr == 10 {
			r[rLength] = uint(i) + 1
			rLength++
		}
	}

	if readonly {
		return CreateUintArray(r)
	} else {
		return r
	}
}
