package fsevents

/*
#cgo LDFLAGS: -framework CoreServices
#include <stdlib.h>
#include <CoreServices/CoreServices.h>
FSEventStreamRef fswatch_stream_for_paths(char** paths, int paths_n);
*/
import "C"
import "unsafe"

var callbackers = make(map[C.FSEventStreamRef]chan []PathEvent)

type PathEvent struct {
}

func WatchPaths(paths []string) chan []PathEvent {
  successChan := make(chan C.FSEventStreamRef)

  go func() {
    var cpaths []*C.char
    for _, dir := range paths {
      path := C.CString(dir)
      defer C.free(unsafe.Pointer(path))
      cpaths = append(cpaths, path)
    }

    stream := C.fswatch_stream_for_paths(&cpaths[0], C.int(len(cpaths)))

    ok := C.FSEventStreamStart(stream) != 0
    if ok {
      successChan <- stream
      C.CFRunLoopRun()
    } else {
      successChan <- nil
    }
  }()

  stream := <-successChan

  if stream == nil {
    return nil
  }

  newChan := make(chan []PathEvent)
  callbackers[stream] = newChan
  return newChan
}

//export watchDirsCallback
func watchDirsCallback(stream C.FSEventStreamRef, count C.size_t, paths **C.char, flags *C.FSEventStreamEventFlags) {
  ch := callbackers[stream]
  ch <- nil

  //for _, ch := range fileSystemChangeObservers {
  //  ch <- []PathEvent{}
  //}
}
