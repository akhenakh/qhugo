package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Global state for MCP
var (
	mcpClient *client.Client
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
)

// InitBackend initializes context
//
//export InitBackend
func InitBackend() {
	ctx, cancel = context.WithCancel(context.Background())
}

// ConnectMCP connects to an MCP server (e.g., qmd)
//
//export ConnectMCP
func ConnectMCP(command *C.char) *C.char {
	cmdStr := C.GoString(command)
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return C.CString(`{"error": "empty command"}`)
	}

	c, err := client.NewStdioMCPClient(parts[0], nil, parts[1:]...)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	mu.Lock()
	if mcpClient != nil {
		mcpClient.Close()
	}
	mcpClient = c
	mu.Unlock()

	connCtx, _ := context.WithTimeout(ctx, 5*time.Second)
	err = c.Start(connCtx)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "failed to start: %s"}`, err.Error()))
	}

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "QtMarkdown", Version: "1.0.0"}

	_, err = c.Initialize(connCtx, initReq)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "failed to init: %s"}`, err.Error()))
	}

	return C.CString(`{"status": "connected"}`)
}

// CallMCPTool executes a tool on the connected MCP server
//
//export CallMCPTool
func CallMCPTool(name *C.char, argsJson *C.char) *C.char {
	mu.Lock()
	defer mu.Unlock()
	if mcpClient == nil {
		return C.CString(`{"error": "not connected"}`)
	}

	toolName := C.GoString(name)
	argsStr := C.GoString(argsJson)

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		return C.CString(`{"error": "invalid arguments json"}`)
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName
	req.Params.Arguments = args

	res, err := mcpClient.CallTool(ctx, req)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	b, _ := json.Marshal(res)
	return C.CString(string(b))
}

// ReadFileContent returns file content
//
//export ReadFileContent
func ReadFileContent(path *C.char) *C.char {
	b, err := os.ReadFile(C.GoString(path))
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(b))
}

// SaveFileContent writes file content
//
//export SaveFileContent
func SaveFileContent(path *C.char, content *C.char) int {
	err := os.WriteFile(C.GoString(path), []byte(C.GoString(content)), 0644)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

// FreeString frees CStrings created by Go
//
//export FreeString
func FreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {}
