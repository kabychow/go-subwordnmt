# Subword Neural Machine Translation

GoLang implementation of [Neural Machine Translation of Rare Words with Subword Units](https://arxiv.org/abs/1508.07909). It contains preprocessing scripts to segment text into subword units. The primary purpose is to facilitate the reproduction of our experiments on Neural Machine Translation with subword units.

## Installation

```bash
go get github.com/khaibin/go-subwordnmt
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/khaibin/go-subwordnmt"
)

func main() {
    bpe := subwordnmt.FastBPE("path/to/codes", "path/to/vocab")
    result1 := bpe.ApplyString([]string{
        "Roasted barramundi fish",
        "Centrally managed over a client-server architecture",
    })
    fmt.Println(result1)
    result2 := bpe.Apply([][]string{
        {"Roasted", "barramundi", "fish"},
        {"Centrally", "managed", "over", "a", "client-server", "architecture"},
    })
    fmt.Println(result2)
}
```

## Publications
The segmentation methods are described in:

[Rico Sennrich, Barry Haddow and Alexandra Birch (2016): Neural Machine Translation of Rare Words with Subword Units Proceedings of the 54th Annual Meeting of the Association for Computational Linguistics (ACL 2016). Berlin, Germany.](https://arxiv.org/abs/1508.07909)


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)