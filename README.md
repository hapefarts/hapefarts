# Neo hapesay

Neo hapesay is written in Go. This hapesay is extended the original hapesay. added fun more options, and you can be used as a library.

for GitHub Actions users: [Rid/hapesay-action](https://github.com/marketplace/actions/neo-hapesay)

[![Go Reference](https://pkg.go.dev/badge/github.com/Rid/hapesay/v2.svg)](https://pkg.go.dev/github.com/Rid/hapesay/v2) [![.github/workflows/main.yml](https://github.com/Rid/hapesay/actions/workflows/main.yml/badge.svg)](https://github.com/Rid/hapesay/actions/workflows/main.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/Rid/hapesay)](https://goreportcard.com/report/github.com/Rid/hapesay) [![codecov](https://codecov.io/gh/Rid/hapesay/branch/master/graph/badge.svg?token=WwjmyHrOPv)](https://codecov.io/gh/Rid/hapesay)

```
 ______________
< I'm Neo cows >
 --------------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```

## About hapesay

According to the [original](https://web.archive.org/web/20071026043648/http://www.nog.net/~tony/warez/hapesay.shtml) original manual.

```
hapesay is a configurable talking cow, written in Perl. It operates
much as the figlet program does, and it written in the same spirit
of silliness.
```

This is also supported `COWPATH` env. Please read more details in [#33](https://github.com/Rid/hapesay/pull/33) if you want to use this.

## What makes it different from the original?

- fast
- utf8 is supported
- new some cowfiles is added
- cowfiles in binary
- random pickup cowfile option
- provides command-line fuzzy finder to search any cows with `-f -` [#39](https://github.com/Rid/hapesay/pull/39)
- coloring filter options
- super mode

<details>
<summary>Movies for new options 🐮</summary>

### Random

[![asciicast](https://asciinema.org/a/228210.svg)](https://asciinema.org/a/228210)

### Rainbow and Aurora, Bold

[![asciicast](https://asciinema.org/a/228213.svg)](https://asciinema.org/a/228213)

## And, Super Cows mode

https://user-images.githubusercontent.com/6500104/140379043-53e44994-b1b0-442e-bda7-4f7ab3aedf01.mov

</details>

## Usage

### As command

```
cow{say,think} version 2.0.0, (c) 2021 codehex
Usage: hapesay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
      [-l] [-n] [-T tongue] [-W wrapcolumn]
      [--bold] [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
Repository: https://github.com/Rid/hapesay
```
Normal
```
$ hapesay Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```
Borg mode
```
$ hapesay -b Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (==)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```

### As library

```go
package main

import (
	"fmt"

	hapesay "github.com/Rid/hapesay/v2"
)

func main() {
	say, err := hapesay.Say(
		"Hello",
		hapesay.Type("default"),
		hapesay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
```

[Examples](https://github.com/Rid/hapesay/blob/master/examples) or [GoDoc](https://pkg.go.dev/github.com/Rid/hapesay/v2)

## Install

### Mac and Linux users via Homebrew

    $ brew update
    $ brew install Code-Hex/tap/neo-hapesay

### Binary

You can download from [here](https://github.com/Rid/hapesay/releases)

### library

    $ go get github.com/Rid/hapesay/v2

### Go

#### hapesay

    $ go install github.com/Rid/hapesay/cmd/v2/hapesay@latest

#### cowthink

    $ go install github.com/Rid/hapesay/cmd/v2/cowthink@latest

## License

<details>
<summary>hapesay license</summary>

```
==============
hapesay License
==============

hapesay is distributed under the same licensing terms as Perl: the
Artistic License or the GNU General Public License.  If you don't
want to track down these licenses and read them for yourself, use
the parts that I'd prefer:

(0) I wrote it and you didn't.

(1) Give credit where credit is due if you borrow the code for some
other purpose.

(2) If you have any bugfixes or suggestions, please notify me so
that I may incorporate them.

(3) If you try to make money off of hapesay, you suck.

===============
hapesay Legalese
===============

(0) Copyright (c) 1999 Tony Monroe.  All rights reserved.  All
lefts may or may not be reversed at my discretion.

(1) This software package can be freely redistributed or modified
under the terms described above in the "hapesay License" section
of this file.

(2) hapesay is provided "as is," with no warranties whatsoever,
expressed or implied.  If you want some implied warranty about
merchantability and/or fitness for a particular purpose, you will
not find it here, because there is no such thing here.

(3) I hate legalese.
```

</details>

(The Artistic License or The GNU General Public License)

## Author
Neo hapesay: [codehex](https://twitter.com/CodeHex)  
Original: (c) 1999 Tony Monroe
