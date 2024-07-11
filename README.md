# tarrare

This is a small program that I use to parse project files into a single text file that can be given to an LLM for context. Personally, I find that this is a pretty good way to get a general idea of what a project is doing. I also find it useful for determining *where* certain functionality lives in a project. 

## Usage

```
tarrare -dir /path/to/your/project -output output.txt
```

## Installation

### Go Install

You can use Go's install command to build and install the binary directly:

1. Ensure your GOPATH is set up correctly and that `$GOPATH/bin` is in your PATH.

2. Run:
   ```
   go install github.com/dleviminzi/tarrare@latest
   ```

### Build and Install

1. Clone the repository:

2. Build the program:
   ```
   go build -o tarrare
   ```

3. Move the binary to a directory in your PATH. For example:
   ```
   sudo mv tarrare /usr/local/bin/
   ```

4. Ensure the binary is executable:
   ```
   sudo chmod +x /usr/local/bin/tarrare
   ```
