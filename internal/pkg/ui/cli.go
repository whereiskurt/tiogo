package ui

import (
	"bufio"
	"fmt"
	"github.com/whereiskurt/tiogo/internal/app"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CLI struct {
	Config      *app.Config
	Workers     *sync.WaitGroup
	WorkerCount int
	Output      *os.File
	Colour      *Colour
	Draw        *Draw
}

func NewCLI(c *app.Config) (cli *CLI) {
	cli = new(CLI)

	cli.Config = c
	cli.Workers = new(sync.WaitGroup)
	cli.WorkerCount, _ = strconv.Atoi(c.WorkerCount)
	cli.Output = c.OutputFileHandle
	cli.Colour = new(Colour)
	cli.Draw = NewDraw(cli)
	if cli.Config.NoColourMode {
		cli.Colour.Disable()
	}

	return
}

func (cli *CLI) Println(line ...interface{}) {
	fmt.Fprintln(cli.Output, fmt.Sprintf("%v", line...))
	return
}

func (cli *CLI) Print(line ...interface{}) {
	fmt.Fprintf(cli.Output, fmt.Sprintf("%v", line...))
	return
}

func (cli *CLI) Stderr(line ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%v", line...))
	return
}

func (cli *CLI) PromptAppConfig() (ok bool) {
	ok = false
	var config = cli.Config

	home := config.HomeFolder

	cli.Println("")
	cli.Println(fmt.Sprintf(cli.Colour.BOLD+"WARN: "+cli.Colour.RESET+"No AppConfiguration file '.tio-cli.yaml' found in homedir '%s' ", home))
	cli.Print(fmt.Sprintf(cli.Colour.BOLD + "Is this your first execution? Need access keys for API usage." + cli.Colour.RESET))
	cli.Println("")
	cli.Println("")
	cli.Draw.Gopher()
	cli.Println("")
	cli.Println("")
	cli.Println(fmt.Sprintf("You must provide the X-ApiKeys '" + cli.Colour.BOLD + "accessKey" + cli.Colour.RESET + "' and '" + cli.Colour.BOLD + "secretKey" + cli.Colour.RESET + "' to the Tenable.IO API."))
	cli.Println(fmt.Sprintf("For complete details see: https://cloud.tenable.com/api#/authorization"))
	cli.Println("")

	reader := bufio.NewReader(os.Stdin)
	cli.Print(fmt.Sprintf("Enter required " + cli.Colour.BOLD + "'accessKey'" + cli.Colour.RESET + ": "))
	config.AccessKey, _ = reader.ReadString('\n')
	config.AccessKey = strings.TrimSpace(config.AccessKey)
	if len(config.AccessKey) != 64 {
		cli.Println(fmt.Sprintf("Invalid accessKey '%s' length %d not 64.\n\n", config.AccessKey, len(config.AccessKey)))
		return
	}

	cli.Print(fmt.Sprintf("Enter required " + cli.Colour.BOLD + "'secretKey'" + cli.Colour.RESET + ": "))
	config.SecretKey, _ = reader.ReadString('\n')
	config.SecretKey = strings.TrimSpace(config.SecretKey)
	if len(config.SecretKey) != 64 {
		cli.Println(fmt.Sprintf("Invalid secretKey '%s' length %d not 64.\n\n", config.SecretKey, len(config.SecretKey)))
		return
	}
	config.CacheKey = fmt.Sprintf("%s%s", config.AccessKey[:16], config.SecretKey[:16])

	cli.Println("")
	cli.Print(fmt.Sprintf("Save AppConfiguration file? [yes or " + cli.Colour.BOLD + "no (default is 'no')" + cli.Colour.RESET + "): "))
	shouldSave, _ := reader.ReadString('\n')
	cli.Println("")

	if len(shouldSave) > 0 && strings.ToUpper(shouldSave)[0] == 'Y' {
		cli.Println(fmt.Sprintf("Creating default '.tio-cli.yaml' in '%s' .", home))

		file, err := os.Create(home + "/.tio-cli.yaml")
		if err != nil {
			cli.Println(fmt.Sprint("Cannot create file:", cli.Colour.BOLD, err, cli.Colour.RESET, "\n\n"))
			return
		}
		defer file.Close()
		cli.Println(fmt.Sprintf("Writing 'accessKey' and 'seretKey'..."))

		fmt.Fprintf(file, "accessKey: %s\n", config.AccessKey)
		fmt.Fprintf(file, "secretKey: %s\n", config.SecretKey)
		fmt.Fprintf(file, "cacheKey: %s%s\n", config.AccessKey[:16], config.SecretKey[:16])
		config.CacheFolder = "./cache"
		fmt.Fprintf(file, "cacheFolder: %s\n", config.CacheFolder)

		t := time.Now()
		ts := fmt.Sprintf("%v", t)
		tzDefault := ts[len(ts)-10:]
		fmt.Fprintf(file, "tzDefault: %s", tzDefault)

		cli.Println(fmt.Sprintf("Done! \nWriting timezone '%v' based on local timezone...", tzDefault))
		cli.Println(fmt.Sprintf("Done! \nSuccessfully created '%v/.tio-cli.yaml'", home))
		cli.Println("")
	}

	ok = true

	return
}
