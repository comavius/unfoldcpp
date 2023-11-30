# unfoldcpp
## purpose
Purpose of unfoldcpp is to unfold cpp-source with complicated inclusion and to return single file.

## usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/comavius/unfoldcpp"
)

func main() {
	single_file, err := unfoldcpp.Unfold("path/to/your/root/file.cpp")
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.Create("path/to/newly/unfolded/file.cpp")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString(single_file)
	if err != nil {
		fmt.Println(err)
	}
}
```

## chenge logs
- 2023-11-30
	Fixed bugs on conditional branch in C++ preprocessing.