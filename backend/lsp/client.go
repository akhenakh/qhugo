package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Client represents an LSP client that communicates with an external LSP server
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	reader *bufio.Reader
	mu     sync.Mutex
	closed bool
	nextID int

	// Callbacks
	onDiagnostic func(uri string, diagnostics []Diagnostic)
	onLog        func(message string)

	// Document management
	documents map[string]*Document
	docMu     sync.RWMutex

	// Request tracking
	pending   map[int]chan *ResponseMessage
	pendingMu sync.Mutex
}

// Document represents an open document
type Document struct {
	URI        string
	Version    int
	Content    string
	LanguageID string
}

// Diagnostic represents a diagnostic message from LSP
type Diagnostic struct {
	Range    Range
	Severity int // 1=Error, 2=Warning, 3=Info, 4=Hint
	Code     interface{}
	Source   string
	Message  string
}

// Range represents a range in a document
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in a document
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Message types for JSON-RPC2
type RequestMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type ResponseMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ErrorObject    `json:"error,omitempty"`
}

type NotificationMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type ErrorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ClientConfig holds configuration for an LSP client
type ClientConfig struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Environment map[string]string `json:"env,omitempty"`
	RootURI     string            `json:"rootUri,omitempty"`
}

// NewClient creates a new LSP client
func NewClient(config ClientConfig, onDiagnostic func(uri string, diagnostics []Diagnostic), onLog func(message string)) (*Client, error) {
	cmd := exec.Command(config.Command, config.Args...)

	// Set environment variables
	if len(config.Environment) > 0 {
		env := cmd.Environ()
		for k, v := range config.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start LSP server: %w", err)
	}

	client := &Client{
		cmd:         cmd,
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		reader:      bufio.NewReader(stdout),
		nextID:      1,
		onDiagnostic: onDiagnostic,
		onLog:       onLog,
		documents:   make(map[string]*Document),
		pending:     make(map[int]chan *ResponseMessage),
	}

	dlog("[LSP Client] Client created, starting readMessages goroutine...")
	// Start message reader and stderr logger
	go client.readMessages()
	go client.logStderr()
	dlog("[LSP Client] Goroutines started")

	return client, nil
}

// Close shuts down the LSP client
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()

	// Send shutdown request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = c.sendRequest(ctx, "shutdown", nil)

	// Send exit notification
	c.sendNotification("exit", nil)

	// Close pipes
	c.stdin.Close()
	c.stdout.Close()
	c.stderr.Close()

	// Kill process if still running
	if c.cmd.Process != nil {
		c.cmd.Process.Kill()
	}

	return c.cmd.Wait()
}

// Initialize sends the initialize request
func (c *Client) Initialize(rootURI string) (*InitializeResult, error) {
	// Ensure rootURI is in proper file:// format
	if rootURI != "" && !strings.HasPrefix(rootURI, "file://") {
		rootURI = "file://" + rootURI
	}

	params := InitializeParams{
		ProcessID:             os.Getpid(),
		ClientInfo:            &ClientInfo{Name: "qhugo", Version: "0.1.0"},
		RootURI:               rootURI,
		InitializationOptions: nil,
		Capabilities: ClientCapabilities{
			TextDocument: &TextDocumentClientCapabilities{
				Synchronization: &TextDocumentSyncClientCapabilities{
					DynamicRegistration: false,
					WillSave:            true,
					WillSaveWaitUntil:   true,
					DidSave:             true,
				},
				Completion: &CompletionClientCapabilities{
					DynamicRegistration: false,
				},
				Hover: &HoverClientCapabilities{
					DynamicRegistration: false,
					ContentFormat:       []string{"markdown", "plaintext"},
				},
				CodeAction: &CodeActionClientCapabilities{
					DynamicRegistration: false,
				},
			},
			Workspace: &WorkspaceClientCapabilities{
				ApplyEdit:              true,
				WorkspaceEdit:          &WorkspaceEditClientCapabilities{DocumentChanges: true},
				DidChangeConfiguration: &DidChangeConfigurationClientCapabilities{DynamicRegistration: false},
			},
		},
	}

	// Only add workspaceFolders if rootURI is set
	if rootURI != "" {
		params.WorkspaceFolders = []WorkspaceFolder{{URI: rootURI, Name: "root"}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.sendRequest(ctx, "initialize", params)
	if err != nil {
		return nil, err
	}

	var initResult InitializeResult
	if err := json.Unmarshal(result, &initResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initialize result: %w", err)
	}

	// Send initialized notification
	c.sendNotification("initialized", InitializedParams{})

	return &initResult, nil
}

// DidOpen sends textDocument/didOpen notification
func (c *Client) DidOpen(uri, languageID, content string) error {
	c.docMu.Lock()
	doc := &Document{
		URI:        uri,
		Version:    1,
		Content:    content,
		LanguageID: languageID,
	}
	c.documents[uri] = doc
	c.docMu.Unlock()

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: languageID,
			Version:    1,
			Text:       content,
		},
	}

	return c.sendNotification("textDocument/didOpen", params)
}

// DidChange sends textDocument/didChange notification with debouncing
func (c *Client) DidChange(uri, content string) error {
	dlog("[LSP Client] DidChange called for %s", uri)
	c.docMu.Lock()
	doc, ok := c.documents[uri]
	if !ok {
		c.docMu.Unlock()
		dlog("[LSP Client] DidChange failed: document not open %s", uri)
		return fmt.Errorf("document not open: %s", uri)
	}

	doc.Version++
	doc.Content = content
	version := doc.Version
	c.docMu.Unlock()
	dlog("[LSP Client] Document %s now at version %d", uri, version)

	// Full document sync for simplicity
	params := DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			URI:     uri,
			Version: version,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: content},
		},
	}

	dlog("[LSP Client] Sending textDocument/didChange notification...")
	err := c.sendNotification("textDocument/didChange", params)
	if err != nil {
		dlog("[LSP Client] Failed to send didChange: %v", err)
	} else {
		dlog("[LSP Client] didChange notification sent successfully")
	}
	return err
}

// DidClose sends textDocument/didClose notification
func (c *Client) DidClose(uri string) error {
	c.docMu.Lock()
	delete(c.documents, uri)
	c.docMu.Unlock()

	params := DidCloseTextDocumentParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
	}

	return c.sendNotification("textDocument/didClose", params)
}

// Hover sends textDocument/hover request
func (c *Client) Hover(uri string, line, character int) (*Hover, error) {
	params := HoverParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{URI: uri},
			Position:     Position{Line: line, Character: character},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.sendRequest(ctx, "textDocument/hover", params)
	if err != nil {
		return nil, err
	}

	var hover Hover
	if err := json.Unmarshal(result, &hover); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hover result: %w", err)
	}

	return &hover, nil
}

// GetDiagnostics returns current diagnostics for a document
func (c *Client) GetDiagnostics(uri string) []Diagnostic {
	// This is populated via the onDiagnostic callback
	// We don't store them here, the UI should maintain its own cache
	return nil
}

// sendRequest sends a JSON-RPC request and waits for response
func (c *Client) sendRequest(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	c.mu.Unlock()

	// Marshal params
	var rawParams json.RawMessage
	if params != nil {
		var err error
		rawParams, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
	}

	req := RequestMessage{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  rawParams,
	}

	// Create response channel
	respChan := make(chan *ResponseMessage, 1)
	c.pendingMu.Lock()
	c.pending[id] = respChan
	c.pendingMu.Unlock()

	// Send request
	if err := c.writeMessage(req); err != nil {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, err
	}

	// Wait for response
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	case <-ctx.Done():
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, ctx.Err()
	}
}

// sendNotification sends a JSON-RPC notification (no response expected)
func (c *Client) sendNotification(method string, params interface{}) error {
	var rawParams json.RawMessage
	if params != nil {
		var err error
		rawParams, err = json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %w", err)
		}
	}

	ntf := NotificationMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  rawParams,
	}

	return c.writeMessage(ntf)
}

// writeMessage writes a JSON-RPC message to the server
func (c *Client) writeMessage(msg interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// LSP message format: Content-Length: <len>\r\n\r\n<json>
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	dlog("[LSP Client] Writing message (%d bytes)", len(data))

	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := c.stdin.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Try to flush if possible
	if flusher, ok := c.stdin.(interface{ Flush() error }); ok {
		if err := flusher.Flush(); err != nil {
			dlog("[LSP Client] Failed to flush: %v", err)
		}
	}

	return nil
}

// readMessages reads and processes messages from the LSP server
func (c *Client) readMessages() {
	dlog("[LSP Client] readMessages goroutine started for client")
	for {
		msg, err := c.readMessage()
		if err != nil {
			if c.closed {
				return
			}
			log.Printf("LSP client read error: %v", err)
			return
		}

		switch m := msg.(type) {
		case *ResponseMessage:
			dlog("[LSP Client] Received response ID=%d", m.ID)
			c.pendingMu.Lock()
			ch, ok := c.pending[m.ID]
			if ok {
				delete(c.pending, m.ID)
			}
			c.pendingMu.Unlock()
			if ok {
				ch <- m
			}

	case *NotificationMessage:
		dlog("[LSP Client] Received notification: %s", m.Method)
		c.handleNotification(m)
		}
	}
}

// readMessage reads a single JSON-RPC message
func (c *Client) readMessage() (interface{}, error) {
	// Read headers
	contentLength := -1
	dlog("[LSP Client] Waiting to read message from server...")
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			dlog("[LSP Client] ReadString error: %v", err)
			return nil, err
		}
		dlog("[LSP Client] Read header line: %q", line)

		line = line[:len(line)-2] // Remove \r\n
		if line == "" {
			break
		}

		if len(line) > 16 && line[:16] == "Content-Length: " {
			fmt.Sscanf(line[16:], "%d", &contentLength)
		}
	}

	if contentLength < 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	// Read content
	data := make([]byte, contentLength)
	if _, err := io.ReadFull(c.reader, data); err != nil {
		return nil, err
	}
	dlog("[LSP Client] Read message content (%d bytes)", contentLength)

	// Parse message
	var base struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Method  string `json:"method"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	if base.Method != "" {
		// It's a notification
		var ntf NotificationMessage
		if err := json.Unmarshal(data, &ntf); err != nil {
			return nil, err
		}
		return &ntf, nil
	}

	// It's a response
	var resp ResponseMessage
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// handleNotification processes server notifications
func (c *Client) handleNotification(ntf *NotificationMessage) {
	switch ntf.Method {
	case "textDocument/publishDiagnostics":
		dlog("[LSP Client] Received publishDiagnostics notification")
		var params PublishDiagnosticsParams
		if err := json.Unmarshal(ntf.Params, &params); err != nil {
			log.Printf("Failed to unmarshal diagnostics: %v", err)
			return
		}
		dlog("[LSP Client] Diagnostics for %s: %d items", params.URI, len(params.Diagnostics))
		if c.onDiagnostic != nil {
			dlog("[LSP Client] Calling onDiagnostic callback")
			c.onDiagnostic(params.URI, params.Diagnostics)
		} else {
			dlog("[LSP Client] onDiagnostic callback is nil!")
		}

	case "window/showMessage":
		var params ShowMessageParams
		if err := json.Unmarshal(ntf.Params, &params); err != nil {
			return
		}
		if c.onLog != nil {
			c.onLog(fmt.Sprintf("[%s] %s", messageTypeString(params.Type), params.Message))
		}

	case "window/logMessage":
		var params LogMessageParams
		if err := json.Unmarshal(ntf.Params, &params); err != nil {
			return
		}
		if c.onLog != nil {
			c.onLog(fmt.Sprintf("[%s] %s", messageTypeString(params.Type), params.Message))
		}
	}
}

// logStderr logs stderr output from the LSP server
func (c *Client) logStderr() {
	scanner := bufio.NewScanner(c.stderr)
	for scanner.Scan() {
		if c.onLog != nil {
			c.onLog(fmt.Sprintf("[LSP stderr] %s", scanner.Text()))
		}
	}
}

func messageTypeString(t int) string {
	switch t {
	case 1:
		return "ERROR"
	case 2:
		return "WARNING"
	case 3:
		return "INFO"
	case 4:
		return "LOG"
	default:
		return "UNKNOWN"
	}
}
