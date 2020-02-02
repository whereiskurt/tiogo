# Welcome to tiogo! [v0.3.2020 :rocket:!!]

[logo]: https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tio gopher wearing red santa hat and bowtie"

![alt text](https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tio gopher wearing red santa hat and bowtie")

# **C**ommand **L**ine **I**nterface to [Tenable.io](https://cloud.tenable.com)

`tiogo` is a command line tool for interacting with the Tenable.io API, written in Go. It follows a general CLI principals of:

```
   ./tio [COMMAND] [SUB-COMMAND] [ACTION] [PARAMS]
```

The `tiogo vm` command currently implements calls to the [Tenable.io Vulnerability API](https://developer.tenable.com/reference) focused on data extracts such as agents, agent-groups, scanners, scans (current/past), vulnerabilities, and assets. Sub-commands such as `export-scans` and `export-assets` make the the `start/status/get` actions easy, requiring minimal parameters. And Sub-commands for `scanners/scans/agents/agent-groups` all default to `list` actions and `--csv` outputs except where `--json` makes more sense. :-)

Tenable offers a variety of for Tenable.io APIs including Web Scanning and Containers. Those APIs may be implmented in the future as `ws` or `container` commands. Today only `vm` exists and is the default and can be ommitted.

Tool written by @whereiskurt and **is not supported or endorsed by Tenable in anyway.**

# Overview

The current primary use case for the `tiogo vm` command is extracting vulns/assets/scans/agents into a SIEM or SOAR system. Because `tiogo` is written in Go it can be complied into a standalone binary for any platform (windows/linux/osx). The binary contains an embeded `jq` binary, the necessary libraries, templates and dependencies to provide a write-once and run-anywhere.

## List your scanners, scan definitions and previous scan run details:

```
  ## Showing `vm` command, it's default and optional.
  $ ./tio vm scanners                  ## Output scanner detail with IP addresses
  $ ./tio vm scans                     ## Output all scans defined

  ## `vm` command not needed, and ommitted
  $ ./tio scanners                     ## Output scanner detail with IP addresses
  $ ./tio scans                        ## Output all scans defined
  $ ./tio scans detail --id=1234       ## Output scan run details for Scan ID '1234'
```

## Agent Group and Agent Lists

```
  $ ./tio agent-groups > agent-groups.20200101.csv
  $ ./tio agents list > agent.list.20200101.csv
```

## Scans Export (JSON/Nessus/CSV/PDF)

Exports of scans/assets/vulnerabilities have a `[start/status/get]` lifecycle. We `export-scans start` our export, then check the `export-scans status` and then `export-scans get` the export. When a scan has run more than once using an `--offset=[0,1,2...]` will get previous results (ie. historical). The default `--offset=0` can be ommited and the current scan will be retrieved.

```
  ###########################
  ## START
  ###########################
  ## Begin export of current scan results in Neuss format (xml)
  $ ./tio export-scans start --id=1234

  ## Begin export of previous scan results in Nessus format (xml)
  $ ./tio export-scans start --offset=1 --id=1234

  ## Begin export of historical (previous previous) scan results in Nessus format (xml)
  $ ./tio export-scans start --offset=2 --id=1234

  ## CSV and PDF [--csv, --pdf]
  $ ./tio export-scans start --id=1234 --csv
  $ ./tio export-scans start --id=1234 --csv --offset=1
  $ ./tio export-scans start --id=1234 --csv --offset=2

  $ ./tio export-scans start --id=1234 --pdf
  $ ./tio export-scans start --id=1234 --pdf --offset=0 --chapter=vuln_by_asset

  $ ./tio export-scans start --id=1234 --pdf --offset=1 --chapter=vuln_by_host
  $ ./tio export-scans start --id=1234 --pdf --offset=2 --chapter=vuln_by_host

  ###########################
  ## STATUS
  ###########################
  $ ./tio export-scans status --id=1234
  $ ./tio export-scans status --id=1234 --csv
  $ ./tio export-scans status --id=1234 --pdf

  $ ./tio export-scans status --id=1234 --offset=1
  $ ./tio export-scans status --id=1234 --csv --offset=1
  $ ./tio export-scans status --id=1234 --pdf --offset=1

  $ ./tio export-scans status --id=1234 --offset=2
  $ ./tio export-scans status --id=1234 --csv --offset=2
  $ ./tio export-scans status --id=1234 --pdf --offset=2

  ###########################
  ## DOWNLOAD
  ###########################
  ## Nessus XML format
  $ ./tio export-scans get --id=1234
  $ ./tio export-scans get --id=1234 --offset=1
  $ ./tio export-scans get --id=1234 --offset=2

  ## JSON (from Nessus XML)
  $ ./tio export-scans query --id=1234 > scan.1234.offset.0.nessus.json

  ## CSV
  $ ./tio export-scans get --id=1234 --csv
  $ ./tio export-scans get --id=1234 --csv --offset=1
  $ ./tio export-scans get --id=1234 --csv --offset=2

  ## PDF
  $ ./tio export-scans get --id=1234 --pdf
  $ ./tio export-scans get --id=1234 --pdf --offset=1
  $ ./tio export-scans get --id=1234 --pdf --offset=2

  ## Convert Nessus XML to JSON using query
  $ ./tio export-scans query --id=1234 > scan.1234.offset.0.nessus.json
  $ ./tio export-scans query --id=1234 --offset=1 > scan.1234.offset.1.nessus.json
  $ ./tio export-scans query --id=1234 --offset=2 > scan.1234.offset.2.nessus.json
```

## Vulnerabilities Export (JSON):

Exports of scans/assets/vulnerabilities have a `[start/status/get]` lifecycle. We `export-vulns start` our export, then check the `export-vulns status` and then `export-vulns get` the export. Using `export-vulns query` allows a `--jqex=<expression>` to be executed on the exported JSON.

```
  $ ./tio export-vulns start --days=365   ## Export 365 days of captured vulns
  $ ./tio export-vulns status             ## Check the status
  ...                                     ##  ... wait until 'FINISHED'
  $ ./tio export-vulns get                ## Download all the chunks
  $ ./tio export-vulns query              ## Dump JSON (--jqex=.)
```

## Assets Export (JSON):

Exports have a `[start/status/get]` lifecycle. We `export-assets start` our export, then check the `export-assets status` and then `export-assets get` the export. Using `export-assets query` allows a `--jqex=<expression>` to be executed on the exported JSON.

```
  $ ./tio export-assets start  ## Start vulns export of a years worth
  $ ./tio export-assets status ## Check the status
  ...                          ##  ... wait until 'FINISHED'
  $ ./tio export-assets get    ## Download chunks
  $ ./tio export-assets query  ## Dump JSON (--jqex=.)
```

# `Dockerfile`

## Using the Dockerfile is a fast way to get 'up and running' if you already have Docker installed and working:

```
$ docker build --tag tiogo:v0.3.2020 .
... [tiogo builds and verbosely outputs]

$ docker run --tty --interactive --rm tiogo:v0.3.2020
root@4f51ab2342123:/tiogo# ./tio help
```

## `> go run cmd\tio.go help`

`tiogo` currently only supports the Vulnerability Management APIs and the defaults to `vm`.

```
root@69e1a9f2bbb2:/tiogo# ./tio help
An interface into the Tenable.io API using Go!

Version v0.3.2020 0521bb94
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
    tio [vm] [SUBCOMMAND] [ACTION ...] [OPTIONS]

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

# Design

## CLI -> Client -> Local Proxy -> Tenable.io

This code is actually three major components CLI/config, Proxy Server and Client:

```
+                                +             +                    +              +
| 1) ./tio.go is called, reads   |             |  2) Start a Proxy  |              |
|    YAML configuration file     |             |       Server       |              |
| +----------------------------+ |             |                    |              |
| |                            | |             | +--------------+   |              |
| |  Command Line Invocation   +----------------->              |   |              |
| | (.tio.yaml configuration)  | | +--------+  | | Proxy Server |   | +----------+ |
| |                            +---> Client +---->              +----->Tenable.io| |
| |                            <---+        <----+              <-----+          | |
| |                            | | +--------+  | +--------------+   | +----------+ |
| +----------------------------+ |             |  4) Relay calls    |              |
|                                |3) Use Client|  to Tenable.io     |              |
|                                |to make calls|                    |              |
|                                |to proxy     |                    |              |
+                                +             +                    +              +

```

I original conceived of this design while working on [tio-cli](https://github.com/whereiskurt/tio-cli/) when Tenable.io backend services were changing frequently and I need a way to insulate my client queries from the Tenable.io responses. Now things are (more) stable and I'm considering no longer maintaining the Proxy Server.

## CLI -> Client -> Tenable.io

You can already acheive the whole 'local proxy' just by pointing the client at `BaseURL` to `cloud.tenable.io` and setting the `DefaultServerStart` to `false` will make the call chain look like this:

```
+                                +             +                    +              +
| 1) ./tio.go is called, reads   |             |                    |              |
|    YAML configuration file     |             |                    |              |
| +----------------------------+ |             |                    |              |
| |                            | |             |                    |              |
| |  Command Line Invocation   +---------------+                    |              |
| | (.tio.yaml configuration)  | | +--------+  |                    | +----------+ |
| |                            +---> Client +------------------------->Tenable.io| |
| |                            <---+        <-------------------------+          | |
| |                            | | +--------+  |                    | +----------+ |
| +----------------------------+ |             |                    |              |
|                                | 2)Use Client|                    |              |
|                                |   to call   |                    |              |
|                                |  Tenable.io |                    |              |
|                                |  API direct |                    |              |
|                                |             |                    |              |
|                                |             |                    |              |
|                                |             |                    |              |
+                                +             +                    +              +

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
