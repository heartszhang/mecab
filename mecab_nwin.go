package mecab

// +build !windows

// #cgo LDFLAGS: -L/usr/local/Cellar/mecab/0.996/lib -lmecab -lstdc++
// #cgo CFLAGS: -I/usr/local/Cellar/mecab/0.996/include
// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
  "unsafe"
  "errors"
)

type MeCab struct {
  _mecab *C.mecab_t
}

func New() (*MeCab, error) {
  _opt := C.CString("-a")
  defer C.free(unsafe.Pointer(_opt))
  _mt := C.mecab_new2(_opt)
  if _mt == nil {
    return nil, MeCabError(nil)
  }
  return &MeCab{_mt}, nil
}

func (this *MeCab) Destroy() {
  C.mecab_destroy(this._mecab)
}

func (this *MeCab) Version() string {
  _ver := C.mecab_version()
  defer C.free(unsafe.Pointer(_ver))
  return C.GoString(_ver)
}

func (this *MeCab) Sparse(x string) (string, error) {
  _str := C.CString(x)
  defer C.free(unsafe.Pointer(_str))
  _v := C.mecab_sparse_tostr(this._mecab, _str)
  if (_v == nil) {
    return "", MeCabError(this)
  }
  return C.GoString(_v), nil
}

func (this *MeCab) NBestSparse(x string, n int) (string, error) {
  _str := C.CString(x)
  defer C.free(unsafe.Pointer(_str))
  //不用释放这块内存
  _v := C.mecab_nbest_sparse_tostr(this._mecab, C.size_t(n), _str)
  if _v == nil {
    return "", MeCabError(this)
  }

  return C.GoString(_v), nil
}

type Node struct {
  Surface  string
  Cost     int64
  IsBest   int
  Position int
  stat     int
}

func (this *MeCab) Sparse2(x string) ([]Node, error) {
  _str := C.CString(x)
  defer C.free(unsafe.Pointer(_str))
  _node := C.mecab_sparse_tonode(this._mecab, _str)

  var nodes []Node
  for current := _node; current != nil; current = current.next {
    node := Node{
      C.GoStringN(current.surface, C.int(current.length)),
      int64(current.cost),
      int(current.isbest),
      int(uintptr(unsafe.Pointer(current.surface)) - uintptr(unsafe.Pointer(_str))),
      int(current.stat),
    }
    //MECAB_NOR_NODE=0/MECAB_UNK_NODE=1/MECAB_BOS_NODE=2/MECAB_EOS_NODE=3
    if (node.stat == 1 || node.stat == 0 ) {
      nodes = append(nodes, node)
    }
  }
  return nodes, nil
}
func MeCabError(this *MeCab) error {
  _err := C.mecab_strerror(this._mecab)
  defer C.free(unsafe.Pointer(_err))
  return errors.New(C.GoString(_err))
}