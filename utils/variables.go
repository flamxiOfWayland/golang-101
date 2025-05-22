package utils

var IString string = ""

var IInt int = 0
var IInt16 int16 = -1
var IInt32 int32 = -22
var IUint uint = 0
var IUint16 uint16 = 1
var IUint32 uint32 = 22

var IFloat32 float32 = 0.0
var IFloat64 float64 = 0.0

var IRune rune = 1

var IComplex64 complex64 = 15 + 2i
var IComplex128 complex128 = 15 + 2i

var IPInt *int = &IInt

var IArrayInt [4]int = [4]int{1, 2, 3, 4}

var ISliceInt []int = []int{1, 2, 3, 4}

var IMapIntString map[int]string = map[int]string{1: "one", 2: "two"}

var IInterface interface{} = "asadsad"

var IChannel1 chan int
var IChannel2 chan<- int
var IChannel3 <-chan int
