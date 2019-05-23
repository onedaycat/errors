package errors

import (
    "fmt"
    "runtime"
    "strings"
)

type Stacktrace struct {
    Frames []*StacktraceFrame `json:"frames"`
}

func (f *Stacktrace) String() string {
    sb := &strings.Builder{}
    for _, frame := range f.Frames {
        sb.WriteString(fmt.Sprintf("%s\t%s:%d\n", frame.Function, frame.AbsolutePath, frame.Lineno))
    }

    return sb.String()
}

func (f *Stacktrace) Strings() []string {
    strs := make([]string, len(f.Frames))
    for i, frame := range f.Frames {
        strs[i] = fmt.Sprintf("%s %s:%d", frame.Function, frame.AbsolutePath, frame.Lineno)
    }

    return strs
}

type StacktraceFrame struct {
    Filename string `json:"filename,omitempty"`
    Function string `json:"function,omitempty"`
    Module   string `json:"module,omitempty"`

    Lineno       int    `json:"lineno,omitempty"`
    Colno        int    `json:"colno,omitempty"`
    AbsolutePath string `json:"abs_path,omitempty"`
}

func NewStacktrace(skip int) *Stacktrace {
    var frames []*StacktraceFrame

    callerPcs := make([]uintptr, 50)
    numCallers := runtime.Callers(skip+2, callerPcs)

    // If there are no callers, the entire stacktrace is nil
    if numCallers == 0 {
        return nil
    }

    callersFrames := runtime.CallersFrames(callerPcs)

    for {
        fr, more := callersFrames.Next()
        frame := &StacktraceFrame{AbsolutePath: fr.File, Filename: fr.File, Lineno: fr.Line}
        frame.Module, frame.Function = functionName(fr.Function)

        // `runtime.goexit` is effectively a placeholder that comes from
        // runtime/asm_amd64.s and is meaningless.
        if frame.Module == "runtime" && frame.Function == "goexit" {
            frame = nil
        }
        if frame != nil {
            frames = append(frames, frame)
        }
        if !more {
            break
        }
    }
    // If there are no frames, the entire stacktrace is nil
    if len(frames) == 0 {
        return nil
    }
    // Optimize the path where there's only 1 frame
    if len(frames) == 1 {
        return &Stacktrace{frames}
    }
    // Sentry wants the frames with the oldest first, so reverse them
    for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
        frames[i], frames[j] = frames[j], frames[i]
    }
    return &Stacktrace{frames}
}

func functionName(fName string) (pack string, name string) {
    name = fName
    // We get this:
    //	runtime/debug.*T·ptrmethod
    // and want this:
    //  pack = runtime/debug
    //	name = *T.ptrmethod
    if idx := strings.LastIndex(name, "."); idx != -1 {
        pack = name[:idx]
        name = name[idx+1:]
    }
    name = strings.Replace(name, "·", ".", -1)
    return
}
