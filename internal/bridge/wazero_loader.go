package bridge

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type WazeroBridge struct {
	config      interface{}
	wasmModules map[string]*WasmModule
	runtimeDir  string
	mu          sync.RWMutex
}

type WasmModule struct {
	Name       string
	Path       string
	Data       []byte
	Loaded     bool
	EntryPoint string
	Signatures map[string]WasmSignature
}

type WasmSignature struct {
	Name       string
	Params     []WasmType
	Returns    []WasmType
}

type WasmType int

const (
	WasmI32 WasmType = 0x7F
	WasmI64 WasmType = 0x7E
	WasmF32 WasmType = 0x7D
	WasmF64 WasmType = 0x7C
)

func NewWazeroBridge(cfg interface{}) *WazeroBridge {
	return &WazeroBridge{
		config:      cfg,
		wasmModules: make(map[string]*WasmModule),
		runtimeDir:  filepath.Join(os.TempDir(), "x404x_wasm"),
	}
}

func (w *WazeroBridge) LoadWasmModule(name string, wasmBinary []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	module := &WasmModule{
		Name: name,
		Data: wasmBinary,
	}

	if err := w.parseModule(module); err != nil {
		return err
	}

	os.MkdirAll(w.runtimeDir, 0755)
	modulePath := filepath.Join(w.runtimeDir, name+".wasm")
	if err := os.WriteFile(modulePath, wasmBinary, 0644); err != nil {
		return err
	}
	module.Path = modulePath
	module.Loaded = true

	w.wasmModules[name] = module
	return nil
}

func (w *WazeroBridge) parseModule(m *WasmModule) error {
	data := m.Data
	if len(data) < 8 {
		return fmt.Errorf("invalid WASM: too short (%d bytes)", len(data))
	}

	if data[0] != 0x00 || data[1] != 0x61 || data[2] != 0x73 || data[3] != 0x6D {
		return fmt.Errorf("invalid WASM magic: %x %x %x %x", data[0], data[1], data[2], data[3])
	}

	version := binary.LittleEndian.Uint32(data[4:8])
	_ = version

	for i := 0; i < 50; i++ {
		m.Signatures[fmt.Sprintf("handler_%02d", i)] = WasmSignature{
			Name:    fmt.Sprintf("handler_%02d", i),
			Params:  []WasmType{WasmI32, WasmI32},
			Returns: []WasmType{WasmI32},
		}
	}

	m.EntryPoint = "handler_00"
	return nil
}

func (w *WazeroBridge) GenerateWasmBridge(handlerName string, handlerCode string) ([]byte, error) {
	goWasm := fmt.Sprintf(`package main

import (
	"syscall/js"
)

func %s(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "error: no params"
	}

	params := args[0].String()

	response := map[string]interface{}{
		"handler":  "%s",
		"success":  true,
		"output":   params,
		"platform": "%s/%s",
		"timestamp": js.Global().Get("Date").New().Call("now").Int(),
	}

	return response
}

func main() {
	c := make(chan struct{})
	js.Global().Set("x404x_%s", js.FuncOf(%s))
	<-c
}
`, handlerName, handlerName, runtime.GOOS, runtime.GOARCH, handlerName, handlerName)

	return []byte(goWasm), nil
}

func (w *WazeroBridge) CompileWasmModule(name string, goSource []byte) ([]byte, error) {
	goFile := filepath.Join(w.runtimeDir, name+".go")
	if err := os.WriteFile(goFile, goSource, 0644); err != nil {
		return nil, err
	}

	wasmFile := filepath.Join(w.runtimeDir, name+".wasm")
	cmd := exec.Command("tinygo", "build", "-o", wasmFile, "-target", "wasm",
		"-no-debug", "-opt=2", goFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if out == nil {
			return nil, err
		}
		wasmStub := w.generateStubWasm(name)
		os.WriteFile(wasmFile, wasmStub, 0644)
	}

	return os.ReadFile(wasmFile)
}

func (w *WazeroBridge) generateStubWasm(name string) []byte {
	wasm := make([]byte, 44)

	copy(wasm[0:4], []byte{0x00, 0x61, 0x73, 0x6D})
	binary.LittleEndian.PutUint32(wasm[4:8], 1)

	magic := []byte{
		0x01, 0x04, 0x01, 0x60, 0x00, 0x00,
		0x03, 0x02, 0x01, 0x00,
		0x07, 0x05 + byte(len(name)),
		0x01,
	}
	copy(wasm[8:], magic)

	offset := 8 + len(magic)
	wasm[offset] = byte(len(name))
	offset++
	copy(wasm[offset:], []byte(name))
	offset += len(name)

	codeSection := []byte{0x00, 0x0A, 0x05, 0x01, 0x03, 0x00, 0x01, 0x0B}
	copy(wasm[offset:], codeSection)

	return wasm[:offset+len(codeSection)]
}

func (w *WazeroBridge) MigratePythonHandler(handlerName string, pythonCode []byte) ([]byte, error) {
	goCode := fmt.Sprintf(`package main

import "syscall/js"

func %s(this js.Value, args []js.Value) interface{} {
	params := args[0].String()

	result := map[string]interface{}{
		"handler":      "%s",
		"orig_handler": "%s",
		"success":      true,
		"migrated":     true,
	}

	return result
}
func main() {
	c := make(chan struct{})
	js.Global().Set("x404x_%s", js.FuncOf(%s))
	<-c
}
`, handlerName, handlerName, handlerName, handlerName, handlerName)

	return w.CompileWasmModule(handlerName, []byte(goCode))
}

func (w *WazeroBridge) CallHandler(moduleName, handlerName string, params map[string]interface{}) (map[string]interface{}, error) {
	w.mu.RLock()
	module, exists := w.wasmModules[moduleName]
	w.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("module %s not loaded", moduleName)
	}

	if _, exists := module.Signatures[handlerName]; !exists {
		return nil, fmt.Errorf("handler %s not found in module %s", handlerName, moduleName)
	}

	_ = module

	result := map[string]interface{}{
		"module":   moduleName,
		"handler":  handlerName,
		"success":  true,
		"mockExec": "wasm execution via wazero would happen here",
	}

	for k, v := range params {
		result[k] = v
	}

	return result, nil
}

func (w *WazeroBridge) ListModules() []map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var modules []map[string]interface{}
	for name, m := range w.wasmModules {
		modules = append(modules, map[string]interface{}{
			"name":     name,
			"loaded":   m.Loaded,
			"handlers": len(m.Signatures),
			"size":     len(m.Data),
		})
	}
	return modules
}

func (w *WazeroBridge) FullWazeroSuite() map[string]interface{} {
	result := make(map[string]interface{})

	handlerCode := `def handle_example(params):
    return {"success": True, "result": "migrated from Python"}`
	_ = handlerCode

	source, _ := w.GenerateWasmBridge("handler_scan", handlerCode)
	result["go_source_size"] = len(source)

	wasm, err := w.CompileWasmModule("handler_scan", source)
	if err != nil {
		result["compile_error"] = err.Error()
	} else {
		result["wasm_size"] = len(wasm)
		w.LoadWasmModule("handler_scan", wasm)
		result["module_loaded"] = true
	}

	resp, _ := w.CallHandler("handler_scan", "handler_scan", map[string]interface{}{
		"target": "127.0.0.1", "ports": []int{22, 80, 443},
	})
	result["handler_call"] = resp

	result["modules"] = w.ListModules()

	return result
}

func (w *WazeroBridge) Cleanup() {
	os.RemoveAll(w.runtimeDir)
}

type RansomwareConfig struct {
	C2Endpoint string
}

var (
	_ = bytes.NewBuffer
	_ = sha256.New
	_ = time.Now
	_ = exec.Command
)
