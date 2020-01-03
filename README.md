# Welcome to tiogo! [v0.2.2019 :rocket:!!]

[logo]: https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tio gopher wearing red santa hat and bowtie"

![alt text](https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tio gopher wearing red santa hat and bowtie")

# **C**ommand **L**ine **I**nterface to [Tenable.io](https://cloud.tenable.com)

`tiogo` is a command line tool for interacting with the Tenable.io API, written in Go. Currently working with the [Tenable.io vulnerability API](https://developer.tenable.com/reference) you can make various calls to extract data for agents, agent-groups, export-vulns, export-assets, export-scans, scanners, and scans.

The tool is written by KPH (@whereiskurt) and **is not supported or endorsed by Tenable in anyway.**

# Overview
The primary use case for `tiogo` is extracting vulns/assets/scans/agents into a SIEM or SOAR system. Because `tiogo` is written in Go it can be complied into a standalone binary for any platform (windows/linux/osx). The binary contains all of the necessary libraries, templates and dependencies to provide a write-once run-anywhere approach.

Using `tiogo` you can can quickly list your scans:
------

```
  $ ./tio scanners                     ## Output scanner detail with IP addresses 
  $ ./tio scans                        ## Output all scans
  $ ./tio scans detail --id=1234       ## Output scan details for Scan ID '1234'
```

Output all of you agent groups and agents
------
```
  $ ./tio agent-groups
  $ ./tio agents list
```

Export a scan and dump as JSON (from xml)
------
```
  $ ./tio export-scans start --id=1234    ## Start an export of scan id 1234
  $ ./tio export-scans status --id=1234   ## Check the status
  ...                                     ##  ... wait until 'READY'
  $ ./tio export-scans get --id=1234      ## Get the created export
  $ ./tio export-scans query --id=1234    ## Dump JSON from Nesus format (XML)
```

Export vulnerabilities as JSON:
------
```
  $ ./tio export-vulns start --days=365   ## Start vulns export of a years worth
  $ ./tio export-vulns status             ## Check the status
  ...                                     ##  ... wait until 'FINISHED'
  $ ./tio export-vulns get                ## Download chunks
  $ ./tio export-vulns query              ## Dump JSON (--jqex=.)
```

# `Dockerfile`

Using the Dockerfile is a fast way to get 'up and running' if you already have Docker installed and working:
------
```
$ docker build --tag tiogo:v0.2.2019 .
... [tiogo builds and verbosely outputs]

$ docker run --tty --interactive --rm tiogo:v0.2.2019
root@4f51ab2342123:/tiogo# ./tio help
```

## `> go run cmd\tio.go help`

`tiogo` currently only supports the Vulnerability Management APIs and the defaults to `vm`.

```
root@69e1a9f2bbb2:/tiogo# ./tio help
An interface into the Tenable.io API using Go!

Version v0.2.2019 71fee112
                 ,_---~~~~~----._
          _,,_,*^____      _____''*g*\"*,
         / __/ /'     ^.  /      \ ^@q   f
        [  @f | @))    |  | @))   l  0 _/
         \'/   \~____ / __ \_____/    \
          |           _l__l_           I
          }          [______]           I
          ]            | | |            |
          ]             ~ ~             |
          |                            |

        [[@https://gist.github.com/belbomemo]]

Find more information at:
    https://github.com/whereiskurt/tiogo/

Usage:
    tio [SUBCOMMAND] [ACTION ...] [OPTIONS]

Sub-commands:
    help, agents, agent-groups, scans, scanners, export-vulns, export-assets, export-scans, cache

VM Options:
    Selection modifiers:
      --id=[unique id]
      --name=[string]
      --regex=[regular expression]
      --jqex=[jq expression]

Output Modes:
      --csv   Set table outputs to comma separated files [ie. good for Excel + Splunk, etc.]
      --json  Set table outputs to JSON [ie. good for integrations and jq manipulations.]

VM Actions and Examples:
    $ tio agents
    $ tio agent-groups
    $ tio scans
    $ tio export-vulns [start|status|get]
    $ tio export-assets [start|status|get]
    $ tio export-scans [start|status|get] --id=123
    $ tio cache clear all
```

## UserHomeDir and `.tiogo/cache/`

When you run `tiogo` for the first time it will ask you for your AccessKey and SecretKey:

```
  root@d173934e91b2:/tiogo# ./tio agent-groups

  WARN: User configuration file '/root/.tiogo.v1.yaml' not found.

  Tenable.io access keys and secret keys are required for all endpoints.
  You must provide X-ApiKeys header 'accessKey' and 'secretKey' values.
  For complete details see: https://developer.tenable.com/

  Enter Tenable.io'AccessKey': df9db9a933d7480be0a902fa1f5df9db9a933d7df9db9a933d7480be02f5ab21
  Enter Tenable.io'SecretKey': 575fd980bc3685d575fd980bc3685d575fd980bc3685d575fd980bc3685df5ab

  Save configuration file? [yes or default:yes]: yes

  Creating default configuration file '/root/.tiogo.v1.yaml' ...
  Done!

  Successfully wrote user configuration file '/root/.tiogo.v1.yaml'.
```

Saving create configuration file in your user's homefolder `.tiogo.v1.yaml` and ultimately create a `.tiogo/cache/` folder hierarchy. The cache folder contains all of the raw JSON returned from Tenable.io under `.tiogo/cache/server/*`

**By default a 'CacheKey' is not set and the results are stored in plaintext.**

## > tio help export-vulns|export-assets

Use `tiogo` you can easily extract all of the vulnerabilities and assets into a collection of files, and query them using built JSON query processor `jq`.

```
root@d173934e91b2:/tiogo# ./tio help export-vulns

```

## Some details about the code:

I've [curated a YouTube playlist](https://www.youtube.com/playlist?list=PLa1qVAzg1FHthbIaRRbLyA4sNE4PmLmn6) of videos which help explain how I ended up with this structure and 'why things are the way they are.' I've leveraged 'best practices' I've seen and that have been explicted called out by others. Of course **THERE ARE SOME WRINKLES** and few **PURELY DEMONSTRATION** portions of code. I hope to be able to keep improving on this.

- [x] Fundamental Go features like tests, generate, templates, go routines, contexts, channels, OS signals, HTTP routing, build/tags, constraints, "ldflags",
- [x] Uses [`cobra`](https://github.com/spf13/cobra) and [`viper`](https://github.com/spf13/viper) (without func inits!!!)
  - Cleanly separated CLI/configuration invocation from client library calls - by calling `viper.Unmarshal` to transfer our `pkg.Config`
  - **NOTE**: A lot of sample Cobra/Viper code rely on `func init()` making it more difficult to reuse.
- [x] Using [`vfsgen`](https://github.com/shurcooL/vfsgen) in to embed templates into binary
  - The `config\template\*` contain all text output and is compiled into a `templates_generate.go` via [`vfsgen`](https://github.com/shurcooL/vfsgen) for the binary build
- [x] Logging from the [`logrus`](https://github.com/sirupsen/logrus) library and written to `log/`
- [x] Cached response folder `.cache/` with entries from the Server, Client and Services
  - The server uses entries in `.cache/` instead of making Tenable.io calls (when present.)
- [x] [Retry](https://github.com/matryer/try) using @matryer's idiomatic `try.Do(..)`
- [x] Instrumentation with [`prometheus`](https://prometheus.io/) in the server and client library
  - [Tutorials](https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-4-instrumenting-code-in-go-and-java/)
- [x] HTTP serving/routing with middleware from [`go-chi`](https://github.com/go-chi/chi)
  - Using `NewStructuredLogger` middleware to decorate each route with log output
  - `ResponseHandler` to pretty print JSON with [`jq`](https://stedolan.github.io/jq/)
  - Custom middlewares (`InitialCtx`,`ExportCtx`) to handle creating Context from HTTP requests
- [x] An example Dockerfile and build recipe `(docs/recipe/)` for a docker workflow
  - Use `docker build --tag tiogo:bulid .` to create a full golang image
  - Use `docker run -it --rm tiogo:build` to work from with the container
