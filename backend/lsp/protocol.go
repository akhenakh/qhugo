package lsp

// Protocol types for LSP communication

// InitializeParams params for initialize request
type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	ClientInfo            *ClientInfo        `json:"clientInfo,omitempty"`
	RootURI               string             `json:"rootUri,omitempty"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	WorkspaceFolders      []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

// ClientInfo information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// WorkspaceFolder workspace folder info
type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// ClientCapabilities capabilities supported by the client
type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
}

// WorkspaceClientCapabilities workspace-related capabilities
type WorkspaceClientCapabilities struct {
	ApplyEdit              bool                                      `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditClientCapabilities          `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`
}

// WorkspaceEditClientCapabilities workspace edit capabilities
type WorkspaceEditClientCapabilities struct {
	DocumentChanges bool `json:"documentChanges,omitempty"`
}

// DidChangeConfigurationClientCapabilities configuration change capabilities
type DidChangeConfigurationClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// TextDocumentClientCapabilities text document capabilities
type TextDocumentClientCapabilities struct {
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	Completion      *CompletionClientCapabilities       `json:"completion,omitempty"`
	Hover           *HoverClientCapabilities            `json:"hover,omitempty"`
	CodeAction      *CodeActionClientCapabilities       `json:"codeAction,omitempty"`
}

// TextDocumentSyncClientCapabilities sync capabilities
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	WillSave          bool `json:"willSave,omitempty"`
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	DidSave           bool `json:"didSave,omitempty"`
}

// CompletionClientCapabilities completion capabilities
type CompletionClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// HoverClientCapabilities hover capabilities
type HoverClientCapabilities struct {
	DynamicRegistration bool     `json:"dynamicRegistration,omitempty"`
	ContentFormat       []string `json:"contentFormat,omitempty"`
}

// CodeActionClientCapabilities code action capabilities
type CodeActionClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// InitializeResult result from initialize request
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerCapabilities capabilities provided by the server
type ServerCapabilities struct {
	TextDocumentSync       *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	CompletionProvider     interface{}              `json:"completionProvider,omitempty"`
	HoverProvider          interface{}              `json:"hoverProvider,omitempty"` // Can be bool or HoverOptions
	DefinitionProvider     interface{}              `json:"definitionProvider,omitempty"` // Can be bool or DefinitionOptions
	CodeActionProvider     interface{}              `json:"codeActionProvider,omitempty"` // Can be bool or CodeActionOptions
	DiagnosticProvider     interface{}              `json:"diagnosticProvider,omitempty"`
}

// TextDocumentSyncOptions sync options
type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose,omitempty"`
	Change    int  `json:"change,omitempty"` // 0=None, 1=Full, 2=Incremental
	WillSave  bool `json:"willSave,omitempty"`
	DidSave   bool `json:"didSave,omitempty"`
}

// ServerInfo information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializedParams params for initialized notification (empty)
type InitializedParams struct{}

// TextDocumentItem represents a text document
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// TextDocumentIdentifier identifies a text document
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier versioned text document identifier
type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

// DidOpenTextDocumentParams params for didOpen notification
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams params for didChange notification
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// TextDocumentContentChangeEvent content change event
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidCloseTextDocumentParams params for didClose notification
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// TextDocumentPositionParams position params
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// HoverParams params for hover request
type HoverParams struct {
	TextDocumentPositionParams `json:",inline"`
}

// Hover hover result
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// MarkupContent markup content for hover
type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// PublishDiagnosticsParams params for publishDiagnostics notification
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     int          `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// ShowMessageParams params for showMessage notification
type ShowMessageParams struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

// LogMessageParams params for logMessage notification
type LogMessageParams struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}
