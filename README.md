# WPDF

WIP

A go package for reading and extracting the text and metadata from them PDF files.

```go
package main

import (
    "fmt"
	
    "github.com/Wafl97/wpdf"
)

func main() {
    pdf, err := wpdf.Open("myfile.pdf")
    if err != nil {
        panic(err)
    }
    // pages use 1 based indexing
    fmt.Println("page 1:", pdf.Page(1).Text())
    // outputs all the text from page 1
}
```

## Resources used

[PDF Standard](https://opensource.adobe.com/dc-acrobat-sdk-docs/standards/pdfstandards/pdf/PDF32000_2008.pdf)
