<p align="center"><a href="#readme"><img src="https://gh.kaos.st/imc.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/imc"><img src="https://travis-ci.com/essentialkaos/imc.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/imc"><img src="https://goreportcard.com/badge/github.com/essentialkaos/imc" alt="GoReportCard" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-imc-master"><img alt="codebeat badge" src="https://codebeat.co/badges/9e4d9881-0c5f-42e1-a775-a3f2de9550df" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" alt="License" /></a>
</p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`imc` (_Icecast Mission Control_) is a simple terminal dashboard for Icecast.

### Screenshots

<p align="center">
  <img src="https://gh.kaos.st/imc.png" alt="imc preview">
</p>

### Installation

#### From source

Before the initial install, allow git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `imc` from scratch, make sure you have a working Go 1.13+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/imc
```

If you want to update `imc` to latest stable release, do:

```
go get -u github.com/essentialkaos/imc
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/imc/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) imc
```

### Usage

```
Usage: imc {options}

Options

  --host, -H host            URL of Icecast instance
  --user, -U username        Admin username
  --password, -P password    Admin password
  --interval, -i seconds     Update interval in seconds (1-600)
  --help, -h                 Show this help message
  --version, -v              Show version

Examples

  imc -H http://192.168.0.1:9922 -U superuser -P MySuppaPass
  Connect to Icecast on 192.168.0.1:9922 with custom user and password

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/imc.svg?branch=master)](https://travis-ci.com/essentialkaos/imc) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/imc.svg?branch=develop)](https://travis-ci.com/essentialkaos/imc) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
