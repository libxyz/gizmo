# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go workspace containing a streaming JSON parser library. The core functionality is a low-memory, character-by-character JSON parser that provides detailed parsing events and JSON path tracking.

## Architecture

**Main Package**: `github.com/Crescent617/x/stream/jsonparser`
- **State Machine Design**: Uses a finite state machine (`State` enum) to process JSON character by character
- **Event System**: Generates detailed events for each character with JSON path tracking (JSON Pointer format)
- **Container Stack**: Tracks nested structures (objects/arrays) using a container stack
- **Memory Efficient**: Processes JSON without loading entire document into memory

**Key Components**:
- `Parser`: Main parser struct with `Push(rune) Event` method
- `Event`: Contains character, event type, and JSON path
- `State`: Enum for parser states (idle, string, number, boolean, null, key)
- `container`: Tracks current object/array context with path building

## Development Commands

Use the Justfile for common tasks:

```bash
just test      # Run all tests across packages
just lint      # Run golangci-lint across packages
```

**Manual Commands**:
```bash
# Run tests for the stream package
cd stream && go test ./...

# Run linter for the stream package
cd stream && golangci-lint run

# Run specific test
cd stream && go test -v ./jsonparser -run TestParser_Simple

# Run all tests with coverage
cd stream && go test -cover ./...
```

## Key Patterns

**Usage Pattern**:
```go
parser := NewParser()
for _, r := range jsonString {
    event := parser.Push(r)
    // Process event based on event.Type and event.Path
}
```

**JSON Path Format**: Uses JSON Pointer format (`$.key[0].nested`)

**Event Types**: 15 event types covering all JSON tokens (string, number, boolean, null, structural characters, escapes, etc.)

## Testing

Tests are comprehensive and include:
- Basic JSON parsing (objects, arrays, primitives)
- Complex nested structures
- Unicode and escape character handling
- Key/value escape sequences
- Path tracking accuracy

Test files: `stream/jsonparser/parser_test.go`