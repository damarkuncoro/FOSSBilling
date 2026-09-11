package plugins

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	luajson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
)

type LuaEngine struct {
	mu      sync.RWMutex
	scripts map[string]*lua.LState
	extDir  string
}

func NewLuaEngine(extDir string) *LuaEngine {
	if extDir == "" {
		extDir = "extensions"
	}
	_ = os.MkdirAll(extDir, 0755)
	return &LuaEngine{
		scripts: make(map[string]*lua.LState),
		extDir:  extDir,
	}
}

func (e *LuaEngine) LoadExtensions() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	files, err := os.ReadDir(e.extDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".lua") {
			path := filepath.Join(e.extDir, f.Name())
			L := lua.NewState()
			e.setupSDK(L)

			if err := L.DoFile(path); err != nil {
				log.Printf("❌ Failed to load Lua extension %s: %v", f.Name(), err)
				L.Close()
				continue
			}
			e.scripts[f.Name()] = L
			log.Printf("🔌 Loaded Dynamic Lua Extension: %s", f.Name())
		}
	}
	return nil
}

func (e *LuaEngine) setupSDK(L *lua.LState) {
	fb := L.NewTable()
	L.SetGlobal("fb", fb)

	L.SetField(fb, "log", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		log.Printf("[Lua Script] %s", msg)
		return 0
	}))
}

func (e *LuaEngine) ExecuteHook(ctx context.Context, hookName string, data any) (any, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	currentData := data
	for _, L := range e.scripts {
		fn := L.GetGlobal(hookName)
		if fn.Type() != lua.LTFunction {
			continue
		}

		luaData := e.goValueToLua(L, currentData)
		err := L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    1,
			Protect: true,
		}, luaData)

		if err != nil {
			log.Printf("❌ Lua Hook %s Error: %v", hookName, err)
			continue
		}

		ret := L.Get(-1)
		L.Pop(1)
		if ret != lua.LNil {
			currentData = e.luaValueToGo(ret)
		}
	}

	return currentData, nil
}

func (e *LuaEngine) goValueToLua(L *lua.LState, v any) lua.LValue {
	b, _ := json.Marshal(v)
	lv, err := luajson.Decode(L, b)
	if err != nil {
		return lua.LNil
	}
	return lv
}

func (e *LuaEngine) luaValueToGo(v lua.LValue) any {
	b, err := luajson.Encode(v)
	if err != nil {
		return nil
	}
	var res any
	_ = json.Unmarshal(b, &res)
	return res
}

func (e *LuaEngine) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, L := range e.scripts {
		L.Close()
	}
}
