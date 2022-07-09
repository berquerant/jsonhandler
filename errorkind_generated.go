// Code generated by "stringer -type ErrorKind -output errorkind_generated.go"; DO NOT EDIT.

package jsonhandler

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Eunknown-0]
	_ = x[EnotJSONRequest-1]
	_ = x[EreadRequestBody-2]
	_ = x[EunmarshalRequestBody-3]
	_ = x[EmarshalResponse-4]
	_ = x[EwriteResponseBody-5]
	_ = x[EhandlerError-6]
}

const _ErrorKind_name = "EunknownEnotJSONRequestEreadRequestBodyEunmarshalRequestBodyEmarshalResponseEwriteResponseBodyEhandlerError"

var _ErrorKind_index = [...]uint8{0, 8, 23, 39, 60, 76, 94, 107}

func (i ErrorKind) String() string {
	if i < 0 || i >= ErrorKind(len(_ErrorKind_index)-1) {
		return "ErrorKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorKind_name[_ErrorKind_index[i]:_ErrorKind_index[i+1]]
}
