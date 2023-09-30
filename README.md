# Pulumi Provider for Jotform

This repository is a "native" Pulumi provider for [Jotform](https://jotform.com). 

It is based on the [Pulumi provider boilerplate](https://github.com/pulumi/pulumi-provider-boilerplate).

*Caveat: This is my first Golang project; I may have missed the mark on Go idioms and standards. Pull requests are welcome.*

## Build
You can build locally, or use the Dockerfile in this repository to build.

### Docker build
Build and run an interactive Docker container

```bash
docker build -t pulumi-jotform-builder .
docker run -it -v $(pwd):/data -w /data --name pulumi-provider-jotform-builder --entrypoint bash pulumi-jotform-builder
```

In teh Docker container's shell, execute `run-jotform-example` (source is in [docker/run-jotform-example](docker/run-jotform-example))

### Local build
To build locally, ensure the following tools are installed and present in your `$PATH`:

* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Go 1.17](https://golang.org/dl/) or 1.latest
* [NodeJS](https://nodejs.org/en/) 14.x.  We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
* [Yarn](https://yarnpkg.com/)
* [TypeScript](https://www.typescriptlang.org/)
* [Python](https://www.python.org/downloads/) (called as `python3`).  For recent versions of MacOS, the system-installed version is fine.
* [.NET](https://dotnet.microsoft.com/download)


### Build the provider and install the plugin

   ```bash
   $ make build install
   ```
   
This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### Test against the example

Create an account on [Jotform](https://jotform.com) and create an API key in [Account settings](https://www.jotform.com/myaccount/api). Set the API key permissions to "Full access". 

Within the docker container you can run `run-jotform-example` (its source is in `docker/run-jotform-example`)

Alternately, replace `YOUR_API_KEY` with the API key in the steps below.

```bash
cd examples/jotform
pulumi stack init jotform-example-dev
pulumi config set --secret jotform_api_key YOUR_JOTFORM_API_KEY 
pulumi up
```

You can view your new form in [My Forms on Jotform](https://www.jotform.com/myforms/).

### Additional Details

This repository depends on the pulumi-go-provider library. For more details on building providers, please check
the [Pulumi Go Provider docs](https://github.com/pulumi/pulumi-go-provider).

## References

Other resources/examples for implementing providers:
* [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
* [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
* [The Easier Way to Create Pulumi Providers in Go](https://www.pulumi.com/blog/pulumi-go-boilerplate-v2/)
