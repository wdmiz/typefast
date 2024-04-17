Command-line typing test, built using [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework

## How to run
*You need to have Go installed (https://go.dev/doc/install)*

### Clone project
Create and enter the directory where project will be cloned:
```
$ mkdir typefast
$ cd typefast
```
Clone `typefast`:
```
$ git clone https://github.com/wdmiz/typefast
```

### Build and run
On Linux or Mac:
```
$ go build
$ ./typefast <flags...>
```

On Windows:
```
$ go build
$ .\typefast.exe <flags...>
```

Alternatively you can run without building (for example, to quickly test changes in code)
### Run without building
Inside directory where project was cloned, run:
```
$ go run . <flags...>
```

### Run flags:
- `-words` how many words should test consist of (default: `100`)
- `-dict` will randomly choose words from this file
- `-text` if specified it will use text from this file as test text (`-words` and `-dict` will be ignored)

## Example run commands
Using `README.md` as dictionary, test length 60 words
```
$ ./typefast -dict README.md -words 60
```
Using `README.md` as test text
```
$ ./typefast -text README.md
```
