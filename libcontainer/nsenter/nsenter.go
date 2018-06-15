// +build linux,!gccgo

package nsenter

/*
#cgo CFLAGS: -Wall
//extern void nsexec();
extern void hyphtest();
void __attribute__((constructor)) init(void) {
	//nsexec();
	hyphtest();
}
*/
import "C"
