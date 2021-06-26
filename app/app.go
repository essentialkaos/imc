package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	ic "pkg.re/essentialkaos/go-icecast.v2"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "imc"
	DESC = "Icecast Mission Control"
	VER  = "1.2.0"
)

const (
	OPT_HOST     = "H:host"
	OPT_USER     = "U:user"
	OPT_PASS     = "P:password"
	OPT_INTERVAL = "i:interval"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_HOST:     {Value: "http://127.0.0.1:8000", Alias: "url"},
	OPT_USER:     {Value: "admin", Alias: "login"},
	OPT_PASS:     {Value: "hackme", Alias: "pass"},
	OPT_INTERVAL: {Value: 15, Min: 1, Max: 600},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

var host string
var icecast *ic.API

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main func
func Init() {
	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError("Options parsing errors:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_VER) {
		showAbout()
		os.Exit(0)
	}

	if options.GetB(OPT_HELP) {
		showUsage()
		os.Exit(0)
	}

	fmtutil.SizeSeparator = " "

	configureIcecastClient()
	checkConnection()

	err := renderGUI()

	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}
}

// configureIcecastClient configures Icecast client
func configureIcecastClient() {
	var err error

	host = options.GetS(OPT_HOST)

	if !strings.HasPrefix(host, "http") {
		host = "http://" + options.GetS(OPT_HOST)
	}

	icecast, err = ic.NewAPI(host, options.GetS(OPT_USER), options.GetS(OPT_PASS))

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

// checkConnection checks connection to Icecast server
func checkConnection() {
	_, err := icecast.ListMounts()

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

// showUsage print usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddOption(OPT_HOST, "URL of Icecast instance", "host")
	info.AddOption(OPT_USER, "Admin username", "username")
	info.AddOption(OPT_PASS, "Admin password", "password")
	info.AddOption(OPT_INTERVAL, "Update interval in seconds {s-}(1-600){!}", "seconds")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-H http://192.168.0.1:9922 -U superuser -P MySuppaPass",
		"Connect to Icecast on 192.168.0.1:9922 with custom user and password",
	)

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(genUsage(), APP))
	case "fish":
		fmt.Printf(fish.Generate(genUsage(), APP))
	case "zsh":
		fmt.Printf(zsh.Generate(genUsage(), optMap, APP))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout shows info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/imc", update.GitHubChecker},
	}

	about.Render()
}
