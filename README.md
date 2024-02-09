[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# ZEN Go

ZEN Engine is a cross-platform, Open-Source Business Rules Engine (BRE). It is written in Rust and provides native 
bindings for NodeJS, Python and Go. ZEN Engine allows to load and execute JSON Decision Model (JDM) from JSON files.

## Installation

```bash
go get github.com/gorules/zen-go
```

## Usage

ZEN Engine is built as embeddable BRE for your **Rust**, **NodeJS**, **Python** or **Go** applications.
It parses JDM from JSON content. It is up to you to obtain the JSON content, e.g. from file system, database or service call.

If you are looking for a complete **BRMS**, take a look at self-hosted [GoRules BRMS](https://gorules.io) or [GoRules Cloud](https://gorules.io).

### Load and Execute Rules

```go
package main

import (
	"fmt"
	"os"
	"path"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("test-data", key)
	return os.ReadFile(filePath)
}

func main() {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose() // Call to avoid leaks

	output, err := engine.Evaluate("rule.json", map[string]any{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(output)
}

```

For more details on rule format and advanced usage, refer to the [main repository](https://github.com/gorules/zen).


