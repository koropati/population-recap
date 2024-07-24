package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/middleware"
	"github.com/koropati/population-recap/routes"
	"github.com/koropati/population-recap/scheduler"
	"github.com/vektra/mockery/mockery"
)

func main() {
	handleCommand()
}

const regexMetadataChars = "\\.+*?()|[]{}^$"

type Config struct {
	fName       string
	fPrint      bool
	fOutput     string
	fOutpkg     string
	fDir        string
	fRecursive  bool
	fAll        bool
	fIP         bool
	fTO         bool
	fCase       string
	fNote       string
	fProfile    string
	fVersion    bool
	quiet       bool
	fkeepTree   bool
	buildTags   string
	fFileName   string
	fStructName string
}

func MyMock() {
	config := parseConfigFromArgs(os.Args)
	handleConfig(config)
}

func handleConfig(config Config) {
	if config.quiet {
		suppressOutput()
	}

	if config.fVersion {
		printVersion()
		return
	}

	if err := validateConfig(config); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	recursive, filter, limitOne := setupFilter(config)

	if config.fkeepTree {
		config.fIP = false
	}

	if config.fProfile != "" {
		createProfile(config.fProfile)
	}

	var osp mockery.OutputStreamProvider
	if config.fPrint {
		osp = &mockery.StdoutStreamProvider{}
	} else {
		osp = &mockery.FileOutputStreamProvider{
			BaseDir:                   config.fOutput,
			InPackage:                 config.fIP,
			TestOnly:                  config.fTO,
			Case:                      config.fCase,
			KeepTree:                  config.fkeepTree,
			KeepTreeOriginalDirectory: config.fDir,
			FileName:                  config.fFileName,
		}
	}

	visitor := &mockery.GeneratorVisitor{
		InPackage:   config.fIP,
		Note:        config.fNote,
		Osp:         osp,
		PackageName: config.fOutpkg,
		StructName:  config.fStructName,
	}

	walker := mockery.Walker{
		BaseDir:   config.fDir,
		Recursive: recursive,
		Filter:    filter,
		LimitOne:  limitOne,
		BuildTags: strings.Split(config.buildTags, " "),
	}

	generated := walker.Walk(visitor)

	if config.fName != "" && !generated {
		fmt.Printf("Unable to find %s in any go files under this path\n", config.fName)
		os.Exit(1)
	}
}

func suppressOutput() {
	os.Stdout = os.NewFile(uintptr(syscall.Stdout), os.DevNull)
}

func printVersion() {
	fmt.Println(mockery.SemVer)
}

func validateConfig(config Config) error {
	if config.fName != "" && config.fAll {
		return errors.New("specify -name or -all, but not both")
	}
	if (config.fFileName != "" || config.fStructName != "") && config.fAll {
		return errors.New("cannot specify -filename or -structname with -all")
	}
	if config.fName != "" && strings.ContainsAny(config.fName, regexMetadataChars) {
		if _, err := regexp.Compile(config.fName); err != nil {
			return fmt.Errorf("invalid regular expression provided to -name: %v", err)
		}
		if config.fFileName != "" || config.fStructName != "" {
			return errors.New("cannot specify -filename or -structname with regex in -name")
		}
	}
	if config.fName == "" && !config.fAll {
		return errors.New("use -name to specify the name of the interface or -all for all interfaces found")
	}
	return nil
}

func setupFilter(config Config) (bool, *regexp.Regexp, bool) {
	var recursive bool
	var filter *regexp.Regexp
	var limitOne bool

	if config.fName != "" {
		recursive = config.fRecursive
		if strings.ContainsAny(config.fName, regexMetadataChars) {
			filter = regexp.MustCompile(config.fName)
		} else {
			limitOne = true
			filter = regexp.MustCompile(fmt.Sprintf("^%s$", config.fName))
		}
	} else if config.fAll {
		recursive = true
		filter = regexp.MustCompile(".*")
	}
	return recursive, filter, limitOne
}

func createProfile(profilePath string) {
	f, err := os.Create(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer pprof.StopCPUProfile()
}

func parseConfigFromArgs(args []string) Config {

	config := Config{}

	flagSet := flag.NewFlagSet(args[1], flag.ExitOnError)

	flagSet.StringVar(&config.fName, "name", "", "name or matching regular expression of interface to generate mock for")
	flagSet.BoolVar(&config.fPrint, "print", false, "print the generated mock to stdout")
	flagSet.StringVar(&config.fOutput, "output", "./mocks", "directory to write mocks to")
	flagSet.StringVar(&config.fOutpkg, "outpkg", "mocks", "name of generated package")
	flagSet.StringVar(&config.fDir, "dir", ".", "directory to search for interfaces")
	flagSet.BoolVar(&config.fRecursive, "recursive", false, "recurse search into sub-directories")
	flagSet.BoolVar(&config.fAll, "all", false, "generates mocks for all found interfaces in all sub-directories")
	flagSet.BoolVar(&config.fIP, "inpkg", false, "generate a mock that goes inside the original package")
	flagSet.BoolVar(&config.fTO, "testonly", false, "generate a mock in a _test.go file")
	flagSet.StringVar(&config.fCase, "case", "camel", "name the mocked file using casing convention [camel, snake, underscore]")
	flagSet.StringVar(&config.fNote, "note", "", "comment to insert into prologue of each generated file")
	flagSet.StringVar(&config.fProfile, "cpuprofile", "", "write cpu profile to file")
	flagSet.BoolVar(&config.fVersion, "version", false, "prints the installed version of mockery")
	flagSet.BoolVar(&config.quiet, "quiet", false, "suppress output to stdout")
	flagSet.BoolVar(&config.fkeepTree, "keeptree", false, "keep the tree structure of the original interface files into a different repository. Must be used with XX")
	flagSet.StringVar(&config.buildTags, "tags", "", "space-separated list of additional build tags to use")
	flagSet.StringVar(&config.fFileName, "filename", "", "name of generated file (only works with -name and no regex)")
	flagSet.StringVar(&config.fStructName, "structname", "", "name of generated struct (only works with -name and no regex)")

	flagSet.Parse(args[2:])

	return config
}

func handleCommand() {
	if len(os.Args) >= 2 {
		switch command := os.Args[1]; command {

		case "server":
			app := bootstrap.NewApp(bootstrap.WithMailer)
			timeout := time.Duration(app.Config.ContextTimeout) * time.Second
			db := app.DB
			defer app.CloseDBConnection()

			// Configure session store (replace with your chosen store)
			store := sessions.NewCookieStore([]byte(app.Config.SessionKey))

			// Use session middleware

			gin := gin.Default()
			gin.Use(middleware.CorsMiddleware())
			gin.Use(sessions.Sessions("mysession", store))

			routeConfig := routes.SetupConfig{
				Config:         app.Config,
				Timeout:        timeout,
				DB:             db,
				CasbinEnforcer: app.CasbinEnforcer,
				Cryptos:        app.Cryptos,
				Gin:            gin,
				Validator:      app.Validator,
				Mailer:         app.Mailer,
			}

			routes.Setup(&routeConfig)
			gin.Run(app.Config.ServerAddress)

		case "scheduler":
			app := bootstrap.NewApp(bootstrap.WithMailer)
			timeout := time.Duration(app.Config.ContextTimeout) * time.Second
			db := app.DB
			defer app.CloseDBConnection()

			cronConfig := scheduler.SetupConfig{
				Config:         app.Config,
				Timeout:        timeout,
				DB:             db,
				CasbinEnforcer: app.CasbinEnforcer,
				Cryptos:        app.Cryptos,
				Mailer:         app.Mailer,
			}

			scheduler.InitCron(&cronConfig)

		case "help":
			log.Printf("Available List Command:\n")
			log.Printf("- go run cmd\\main.go server    (to start server process)\n")
			log.Printf("- go run cmd\\main.go consumer (to start scheduler consumer process)\n")
			log.Printf("- go run cmd\\main.go publisher (to start scheduler publisher process)\n")
		case "mockery":
			MyMock()
		default:
			log.Printf("It's Working!\n")
		}
	} else {
		log.Printf("Program It's Working!, you must select operation to start a session.\n")
		log.Printf("List Command:\n")
		log.Printf("- go run cmd\\main.go server    (to start server process)\n")
		log.Printf("- go run cmd\\main.go consumer (to start scheduler consumer process)\n")
		log.Printf("- go run cmd\\main.go publisher (to start scheduler publisher process)\n")
		log.Printf("- go run cmd\\main.go help      (to see list of command)\n")
	}
}
