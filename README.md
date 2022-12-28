# iken ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/bir/iken/build.yml?branch=master) [![codecov](https://codecov.io/gh/bir/iken/branch/master/graph/badge.svg)](https://codecov.io/gh/bir/iken) [![PkgGoDev](https://pkg.go.dev/badge/github.com/bir/iken)](https://pkg.go.dev/github.com/bir/iken) [![Report card](https://goreportcard.com/badge/github.com/bir/iken)](https://goreportcard.com/report/github.com/bir/iken) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbir%2Fiken.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbir%2Fiken?ref=badge_shield)

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/open-source.svg)](https://forthebadge.com)

**iken** is an _opinionated_ library for building apps in go.

# High Level Opinions

1. Errors should be managed
1. Panics should be managed
1. Developer Experience is critical to adoption
1. Obfuscate nothing - code should be clear to read and trace.
1. Testability is important
1. Dependency injection is critical to testability
1. Global vars are evil
1. Favor easy codegen over obfuscating libraries when feasible

# Concrete Opinions

1. net/http is the preferred handler
2. Postgres is the preferred SQL DB
3. [pgx](https://github.com/jackc/pgx) is the preferred Postgres library
4. [zerolog](https://github.com/rs/zerolog) is the preferred Logger

_Preferred_ is the keyword, as the packages in **iken** ease the support for these libraries, but do not prevent any
other options.

# Motivation

We had several

# Current projects

## errs

Support for cause chaining with a nil check. The excellent pkg.errors does not handle the case where `Cause()` returns
nil.

`WithStack` provides an easy stack traced error with options to ignore depth. Useful for tracking panics caught in
middleware. It also provides some utilities for marshalling to logging for easy of logging.

## httputil

Collection of minor tools for use with HTTP.

### ErrorHandler

Standardized handling of errors in an HTTP request flow.

## pgxzero

## validation


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbir%2Fiken.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbir%2Fiken?ref=badge_large)