# fourbyte [![Go Report Card](https://goreportcard.com/badge/github.com/0xhyim/fourbyte?)](https://goreportcard.com/report/github.com/0xhyim/fourbyte) [![GoDoc](https://pkg.go.dev/badge/github.com/0xhyim/fourbyte.svg)](https://pkg.go.dev/github.com/0xhyim/fourbyte)

> [4byte.directory](https://www.4byte.directory/) API Client for [Go](https://golang.org/).

Check out the endpoints or usage section for more information.

## Documentation
[https://pkg.go.dev/github.com/0xhyim/fourbyte](https://pkg.go.dev/github.com/0xhyim/fourbyte)

## Install

```bash
go get -u github.com/0xhyim/fourbyte
```

## Endpoints

|               Endpoint                |                                Function                                  |
| :-----------------------------------: | :----------------------------------------------------------------------: |
|          /api/v1/signatures           | GetFunctionSignatures, GetFunctionSignatureById, CreateFunctionSignature |
|       /api/v1/event-signatures        |     GetEventSignatures, GetEventSignatureById, CreateEventSignature      |
|       /api/v1/import-solidity         |                   ImportFromSourceCode, ImportFromABI                    |

Refer to [4byte.directory Official API Documentation](https://www.4byte.directory/docs/) for more information.

## Usage

The `fourbyte` package provides access to all the endpoints listed above.

### Basic Function Signature Lookup

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0xhyim/fourbyte"
)

func main() {
	response, err := fourbyte.GetFunctionSignatures(context.Background(), fourbyte.WithHexSignature("0x1249c58b"))
	if err != nil {
		log.Fatal(err)
	}
	// "mint()"
	fmt.Println(response.Signatures[0].TextSignature)
}
```

### Signature Lookup (following Next)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0xhyim/fourbyte"
)

func main() {
	signatures := []fourbyte.Signature{}
	response, err := fourbyte.GetFunctionSignatures(context.Background(), fourbyte.WithTextSignature("WETH"))
	if err != nil {
		log.Fatal(err)
	}
	signatures = append(signatures, response.Signatures...)
	for response.Next != nil {
		response, err = response.Next.Follow(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		signatures = append(signatures, response.Signatures...)
	}
	fmt.Println(signatures)
}
```

## TODOS

- Add tests to project
- Add documentation/comments to code

## License

[MIT License](LICENSE)