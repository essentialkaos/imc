package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/support"
	"github.com/essentialkaos/ek/v12/support/deps"
	"github.com/essentialkaos/ek/v12/support/pkgs"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	ic "github.com/essentialkaos/go-icecast/v2"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "imc"
	DESC = "Icecast Mission Control"
	VER  = "1.2.2"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_HOST     = "H:host"
	OPT_USER     = "U:user"
	OPT_PASS     = "P:password"
	OPT_INTERVAL = "i:interval"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_HOST:     {Value: "http://127.0.0.1:8000", Alias: "url"},
	OPT_USER:     {Value: "admin", Alias: "login"},
	OPT_PASS:     {Value: "hackme", Alias: "pass"},
	OPT_INTERVAL: {Value: 15, Min: 1, Max: 600},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// colorTagApp contains color tag for app name
var colorTagApp string

// colorTagVer contains color tag for app version
var colorTagVer string

// client is Icecast API client
var client *ic.API

// host is Icecast host
var host string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main application function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError("Options parsing errors:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithPackages(pkgs.Collect("icecast,icecast2,icecast-kh")).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	configureIcecastClient()
	checkConnection()

	err := renderGUI()

	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	fmtutil.SizeSeparator = " "

	switch {
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#27}", "{#27}"
	default:
		colorTagApp, colorTagVer = "{*}{b}", "{b}"
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// configureIcecastClient configures Icecast client
func configureIcecastClient() {
	var err error

	host = options.GetS(OPT_HOST)

	if !strings.HasPrefix(host, "http") {
		host = "http://" + options.GetS(OPT_HOST)
	}

	client, err = ic.NewAPI(host, options.GetS(OPT_USER), options.GetS(OPT_PASS))

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

// checkConnection checks connection to Icecast server
func checkConnection() {
	_, err := client.ListMounts()

	if err != nil {
		printError("Can't connect to Icecast server on %s", options.GetS(OPT_HOST))
		printError("Check URL, username and password")
		os.Exit(1)
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(genUsage(), "imc"))
	case "fish":
		fmt.Print(fish.Generate(genUsage(), "imc"))
	case "zsh":
		fmt.Print(zsh.Generate(genUsage(), optMap, "imc"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddOption(OPT_HOST, "URL of Icecast instance", "host")
	info.AddOption(OPT_USER, "Admin username", "username")
	info.AddOption(OPT_PASS, "Admin password", "password")
	info.AddOption(OPT_INTERVAL, "Update interval in seconds {s-}(1-600){!}", "seconds")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-H http://192.168.0.1:9922 -U superuser -P MySuppaPass",
		"Connect to Icecast on 192.168.0.1:9922 with custom user and password",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{"essentialkaos/imc", update.GitHubChecker}
	}

	return about
}
