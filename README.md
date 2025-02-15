# neatfile

Looking for a tool to remove all those comments and unwanted blank lines in a file especially your code, Then `neatfile` is for you.

lets declutter your file and make it neat.

It removes all the commented lines ( except inline comments) and empty lines from the file and prints the output to the standard output or to a file.

![neatfile](demo-short.gif)


## Usage

```bash

# To print the output to the standard output
neatfile <input-file> 

# To print the output to a file
neatfile -o <output-file> <input-file> 

```

If no output file is specified, the output will be printed to the standard output.

## How to install

### From Git Repository
```
git clone https://github.com/AKSarav/neatfile.git
cd NeatFile
make install
```

### Using Go Install

You can alternatively use go install to install the tool.

```
go install github.com/AKSarav/neatfile
```

### Install with Brew

```
brew tap AKSarav/neatfile
brew install neatfile
```

## Contribute

Feel free to contribute to the project and make it better. create your own branch and raise a PR.

## Leave a star if you like the project :star:



