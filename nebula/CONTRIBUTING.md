# Contributing

Clone the repo and cd:
```bash
$ git clone https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula.git
$ cd nebula
```

We use a module from a from our Github. To tell Go that it needs to check our repo, that it is a private one:
```bash
go env -w "GOPRIVATE=github.vodafone.com/*"
```

To run the tests:
```bash
$ make test
```

To run the linter:
```bash
$ make check
```

Run `make help` for a full description.


## Suggested dev tools

Golang official packages.

## Pull Requests

Check our Confluence page on how to create and review [Pull Requests](https://confluence.sp.vodafone.com/display/NAAP/Pull+Requests).

## Rules

Golang is a very complete language that includes several interesting tools by default, such go test, go doc, etc. It is designed to provide full development functionalities from its standard library and its standard tools, thus, if any doubt on how to code or proceed, please refer to the official Go documentation: https://go.dev/doc/effective_go.
Caveats

There are some standards in Go that might mislead or reduce the quality of the code if used without the needed precautions:

    Variable names: Golang suggests using short variable names, as short as one letter. Sometimes this is good, for example in request handlers: `func(w http.ResponseWriter, r *http.Request)`. However, as a standard, please use variable names that describe what the variable is storing.
    Line length: Golang does not define a limit for the length of the lines and it cannot be enforced with the tools provided by them. Be sensible and avoid very large lines that would reduce or complicate the readability of the code.

## Go Tools

Golang provides a wide variety of tools to ease development. If you use VS Code, it will suggest the installation and automatically configure everything, such as go fmt (the tool that will format your code to a certain standard). Please, configure your Go environment to help all of us!

## Code Structure

The code is structured in different modules covering diverse functionalities. Each module is actually a directory under the root project directory, following Golang's basic repository structure. Tests for each file/module should be stored in the same directory as the file.

## Middleware

Nebula provides a 'chain' of middleware functions that utilise the modules mentioned above to perform validations, filtering and security checks for all incoming requests. The wiki has a page on [middleware functions](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/wiki/Middleware-Functions) describing this in more detail.

## Configuration

To customise some of its middleware, Nebula uses a YAML configuration that allows a user to enable certain features, specify a target host and set rate limits where needed. When ran locally on the development machine, this would come in the form of a config.yaml file in the root project directory. The wiki has a page on [Configuration](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/wiki/Configuration) describing YAML configurations in more detail.