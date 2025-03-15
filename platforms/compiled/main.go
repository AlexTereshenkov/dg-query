package main

// make sure to place the `import "C"` directive as a separate import line

/*
#include <stdlib.h>
#include <unistd.h>

char* get_hostname() {
    char* hostname = malloc(256);
    if (hostname == NULL) {
        return NULL;
    }
    if (gethostname(hostname, 256) != 0) {
        free(hostname);
        return NULL;
    }
    return hostname;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	cstr := C.get_hostname()
	if cstr == nil {
		fmt.Println("Failed to get hostname")
		return
	}
	gostr := C.GoString(cstr)
	fmt.Println("Hostname:", gostr)
	C.free(unsafe.Pointer(cstr))
}
