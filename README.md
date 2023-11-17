# go-climatiq
A GO SDK for the Climatiq API

The **go-climatiq** SDK is a Go (Golang) library that facilitates easy integration with the [Climatiq API](https://www.climatiq.io/). Emission factors are crucial in estimating carbon emissions for various activities, industries, or processes. This SDK simplifies the process of fetching emission factors and incorporating them into your applications for accurate carbon footprint calculations.

## Installation

To install the go-climatiq SDK, use the following `go get` command:

```bash
go get github.com/re-cinq/go-climatiq/v2/climatiq
```


## Usage
To use the SDK in your Go application, import the package and create a client instance:

```go
import (
	"github.com/re-cinq/go-climatiq/v2/climatiq"
)

func main() {
    // Create a new climatiq client with the auth token option
    cli := climatiq.NewClient(
        climatiq.WithAuthToken("YOUR_API_KEY")
    )
}
```
Make sure to replace `YOUR_API_KEY` with the actual API key provided by [Climatiq](https://www.climatiq.io/docs/api-reference/authentication)

Currently, the client only supports the climatiq [Search](https://www.climatiq.io/docs/api-reference/search) query. A small code snipit on how to use that can be found in the `example` directory.

## Examples
Check the [example](https://github.com/re-cinq/go-climatiq/tree/main/example) directory for sample code snippets demonstrating how to use the SDK.

## Contributing
If you find any issues or have suggestions for improvements, please open an issue or create a pull request on the GitHub repository. We are always open to contributions and help!

## License
This SDK is licensed under the MIT License - see the LICENSE file for details.

## Contact
For any inquiries or support, please contact the maintainers: <br>
gabi@re-cinq.com brendan@re-cinq.com sebastian@re-cinq.com

Thank you for using the go-climatiq SDK

