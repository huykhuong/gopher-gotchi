package window

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static void moveWindowBy(void* nswin, double dx, double dy) {
	NSWindow* win = (__bridge NSWindow*)nswin;
	NSRect frame = [win frame];
	frame.origin.x += dx;
	frame.origin.y -= dy; // Cocoa Y is flipped vs screen coords
	[win setFrame:frame display:YES animate:NO];
}

static void positionWindowBottomRight(void* nswin, int width, int height, int margin) {
	NSWindow* win = (__bridge NSWindow*)nswin;
	NSScreen* screen = [NSScreen mainScreen];
	NSRect screen_frame = [screen visibleFrame];
	NSRect frame;
	frame.size.width = width;
	frame.size.height = height;
	frame.origin.x = screen_frame.origin.x + screen_frame.size.width - width - margin;
	frame.origin.y = screen_frame.origin.y + margin;
	[win setFrame:frame display:YES animate:NO];
}

static void makeWindowFloating(void* nswin) {
	NSWindow* win = (__bridge NSWindow*)nswin;
	[win setLevel:NSFloatingWindowLevel];
	[win setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
		NSWindowCollectionBehaviorStationary |
		NSWindowCollectionBehaviorIgnoresCycle];
	[win setBackgroundColor:[NSColor clearColor]];
	[win setOpaque:NO];
	[win setHasShadow:YES];
	[win setMovable:NO]; // We handle dragging ourselves via JS
}
*/
import "C"

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	webview "github.com/webview/webview_go"
)

//go:embed index.html
var indexHTML string

//go:embed spritesheet.webp
var Spritesheet string

const (
	windowWidth  = 200
	windowHeight = 340
	windowMargin = 20
)

// PetState is the data sent to the webview UI on each tick.
type PetState struct {
	Level    int      `json:"level"`
	Hunger   int      `json:"hunger"`
	Mood     string   `json:"mood"`
	Message  string   `json:"message"`
	CPULoad  int      `json:"cpuLoad"`
}

// Window wraps a webview window that displays the floating pet UI.
type Window struct {
	wv     webview.WebView
	onQuit func()
}

// New creates a new floating pet window. onQuit is called when the user
// clicks the Quit button inside the webview.
//
// Must be called from the main OS thread (runtime.LockOSThread is applied
// automatically by the webview package's init()).
func New(onQuit func(), dev bool) *Window {
	runtime.LockOSThread()

	w := &Window{
		wv:     webview.New(dev),
		onQuit: onQuit,
	}

	w.wv.SetTitle("Diana")
	w.wv.SetSize(windowWidth, windowHeight, webview.HintFixed)

	if dev {
		cwd, _ := os.Getwd()
		htmlPath := filepath.Join(cwd, "internal", "window", "index.html")
		w.wv.Navigate("file://" + htmlPath)
	} else {
		w.wv.SetHtml(indexHTML)
	}

	// Position and style the native NSWindow after the webview is created.
	w.wv.Dispatch(func() {
		nswin := w.wv.Window()
		C.makeWindowFloating(nswin)
		C.positionWindowBottomRight(
			nswin,
			C.int(windowWidth),
			C.int(windowHeight),
			C.int(windowMargin),
		)
	})

	// JS binding: moveWindowBy(dx, dy) – called by the drag handler in index.html
	w.wv.Bind("moveWindowBy", func(dx, dy float64) {
		w.wv.Dispatch(func() {
			nswin := w.wv.Window()
			C.moveWindowBy(nswin, C.double(dx), C.double(dy))
		})
	})

	// JS binding: quitApp() – called by the Quit button in index.html
	w.wv.Bind("quitApp", func() {
		w.wv.Dispatch(func() {
			w.wv.Terminate()
		})
		if w.onQuit != nil {
			go w.onQuit()
		}
	})

	return w
}

// Update pushes fresh pet state to the webview UI by calling the JS
// updateState() function defined in index.html.
func (w *Window) Update(state PetState) {
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("window: marshal error: %v\n", err)
		return
	}

	js := fmt.Sprintf("if(typeof updateState==='function'){updateState(%s);}", string(data))
	w.wv.Dispatch(func() {
		w.wv.Eval(js)
	})
}

// Run starts the webview event loop. This call blocks until the window is
// closed. It must be called from the main OS thread.
func (w *Window) Run() {
	w.wv.Run()
	w.wv.Destroy()
}

// Terminate closes the window programmatically (e.g. from a signal handler).
func (w *Window) Terminate() {
	w.wv.Terminate()
}

// ensure unsafe is used (satisfies the import)
var _ = unsafe.Pointer(nil)
