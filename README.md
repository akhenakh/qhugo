# QHugo

A tool specialized to manage and edit Hugo blogs natively.

## Features
- Embedded live preview using Hugo's built-in webserver.
- Drag and drop image processing (automatically resizes and copies to `static/img`).
- Native post creation with automatic frontmatter generation.
- Markdown editor with syntax highlighting.

## Prerequisites
- **Hugo**: Ensure `hugo` is installed and available in your system's PATH.

## Build
```bash
git submodule update --init --recursive
cmake -b build
cd build && make
```



## TODO
- LSP support
