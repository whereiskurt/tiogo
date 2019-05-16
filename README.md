# Welcome to tiogo!!
### A modern CLI for Tenable.io written in Go - v0.1 [20190515] :rocket:

[logo]: https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tiogopher"
![alt text](https://github.com/whereiskurt/tiogo/blob/master/docs/images/tiogo.logo.small.png "tiogopher")

## A **C**ommand **L**ine **I**nterface to Tenable.io API 
`tiogo` is a command line tool for interacting with the Tenable.io API, written in golang. It currently only supports a small set of the [Tenable.io vulnerability API](https://developer.tenable.com/reference) around agents, agent-groups, export-vuls and export-assets. 

The tool is written by KPH (@whereiskurt) and **is not supported or endorsed by Tenable in anyway.**

## Overview 
[Tenable.io](https://cloud.tenable.com) is a modern webapp rendered in web browser - aka **G**raphical **U**ser **I**nterfaces (**GUI**).

Alternatively, `tiogo` is a **C**ommand **L**ine **I**nterface (**CLI**) to interact with Tenable.io. Because `tiogo` is written in Go it can be complied into a standalone binary for any platform. The binary contains all of the necessary libraries and dependencies included and provides a write-once run-anywhere approach.

## `Dockerfile`
Using the Dockerfile is a fast way to get 'up and running' if you already have Docker installed and working:
```
$ docker build --tag tiogo:v0.1 .
... [tiogo builds and verbosely outputs]

$ docker run --it --rm tiogo:v0.1
  
root@4f51ab2342123:/tiogo# ./tio help
```
## `> go run cmd\tio.go help`
`tiogo` currently only supports the Vulnerability Management APIs and the defaults to `vm`.

```
root@d173934e91b2:/tiogo# ./tio help

An interface into the Tenable.io API using Go!

Version v0.1.0 132471e4
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
    tio [COMMAND] [SUBCOMMAND] [ACTION ...] [OPTIONS]

Commands:
    vm       Commands for Tenable.io Vulnerability Management [default, can be omitted]
    server   Commands for local proxy and HTTP server instance

Sub-command:
    vm:
      help, scanners, agents, agent-groups, scans, export-vulns

    proxy:
      start, stop

Global Options:
    Verbosity:
      --silent,  -s     Set logging/output level [level1]
      --quiet,   -q     Set logging/output level [level2]
      --info,    -v     Set logging/output level [level3-default]
      --debug,          Set logging/output level [level4]
      --trace,          Output to STDOUT and to log file [level5]
      --level=3         Sets the output verbosity level numerically [default]

For more help:
    $ tio help scanners
    $ tio help agents
    $ tio help agent-groups
    $ tio help export-vulns
    $ tio help scans
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

Bulk Exports of Vulnerabilities
https://developer.tenable.com/reference#exports

Usage:
    tio vm export-vulns [ACTION ...] [OPTIONS]

Action:
    start, status, get, query

Export Vulns Options:
    Selection modifiers:
    --uuid=[unique id]
    --jqex=[jq expression]
    --chunk=[chunk to get, defaults: ALL]
    --critical, --high, --medium, --info  [severity to match for vulnerability]
    --before=[YYYY-MM-DD HH:MM:SS +/-0000 TZ], --after=[YYYY-MM-DD HH:MM:SS +/-0000 TZ] [date boundaries]
    --days=[number of days to bound query to]

Output Modes:
    --json  Set table outputs to JSON [ie. good for integrations and jq manipulations.]

Examples:
    $ tio export-vulns start
    $ tio export-vulns start --after="2019-01-01" --critical
    $ tio export-vulns start --after="2019-01-01 00:00:00 -0400 EDT"
    $ tio export-vulns start --before=="2019-01-31" --critical --high
    $ tio export-vulns start --before="2019-01-31" --days=31 --critical --high
    $ tio export-vulns start --after=2019-01-01 --days=31

    $ tio export-vulns status
    $ tio export-vulns get
    $ tio export-vulns query --jqex="[.asset.ipv4, .asset.operating_system[0]]"
```

## Some details about the code:
- [x] Fundamental Go features like tests, generate, templates, go routines, contexts, channels, OS signals, HTTP routing, build/tags, constraints, "ldflags", 
- [x] Uses [`cobra`](https://github.com/spf13/cobra) and [`viper`](https://github.com/spf13/viper) (without func inits!!!)
  - Cleanly separated CLI/configuration invocation from client library calls - by calling `viper.Unmarshal` to transfer our `pkg.Config`
  - **NOTE**: A lot of sample Cobra/Viper code rely on `func init()` making it more difficult to reuse. 
- [x] Using [`vfsgen`](https://github.com/shurcooL/vfsgen) in to embed templates into binary
    - The `config\template\*` contain all text output and is compiled into a `templates_generate.go` via [`vfsgen`](https://github.com/shurcooL/vfsgen) for the binary build
- [X] Logging from the [`logrus`](https://github.com/sirupsen/logrus) library and written to `log/`
- [x] Cached response folder `.cache/` with entries from the Server, Client and Services
  - The server uses entries in `.cache/` instead of making Tenable.io calls (when present.)
- [x] [Retry](https://github.com/matryer/try) using @matryer's idiomatic `try.Do(..)`
- [X] Instrumentation with [`prometheus`](https://prometheus.io/) in the server and client library
  - [Tutorials](https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-4-instrumenting-code-in-go-and-java/)
- [X] HTTP serving/routing with middleware from [`go-chi`](https://github.com/go-chi/chi)
    - Using `NewStructuredLogger` middleware to decorate each route with log output
    - `ResponseHandler` to pretty print JSON with [`jq`](https://stedolan.github.io/jq/)
    - Custom middlewares (`InitialCtx`,`ExportCtx`) to handle creating Context from HTTP requests
- [x] An example Dockerfile and build recipe `(docs/recipe/)` for a docker workflow
  - Use `docker build --tag tiogo:bulid .` to create a full golang image
  - Use `docker run -it --rm tiogo:build` to work from with the container

I've [curated a YouTube playlist](https://www.youtube.com/playlist?list=PLa1qVAzg1FHthbIaRRbLyA4sNE4PmLmn6) of videos which help explain how I ended up with this structure and 'why things are the way they are.' I've leveraged 'best practices' I've seen and that have been explicted called out by others. Of course **THERE ARE SOME WRINKLES** and few **PURELY DEMONSTRATION** portions of code. I hope to be able to keep improving on this.
