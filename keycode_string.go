// Code generated by "stringer -type KeyCode -trimprefix=Key"; DO NOT EDIT.

package smclcd

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KeyUp-0]
	_ = x[KeyRight-1]
	_ = x[KeyLeft-2]
	_ = x[KeyDown-3]
	_ = x[KeyEnter-4]
	_ = x[KeyCancel-5]
}

const _KeyCode_name = "UpRightLeftDownEnterCancel"

var _KeyCode_index = [...]uint8{0, 2, 7, 11, 15, 20, 26}

func (i KeyCode) String() string {
	if i >= KeyCode(len(_KeyCode_index)-1) {
		return "KeyCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _KeyCode_name[_KeyCode_index[i]:_KeyCode_index[i+1]]
}
