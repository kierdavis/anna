**PLEASE NOTE: I no longer actively maintain this project. Pull requests are welcome, but development has halted and issues are unlikely to gain a response.**

## Installation

This package utilises vendoring (as per the Go 1.5 standard) through the use of
Git submodules; it appears that currently `go get` is not able to install these
packages correctly.

### Install SDL 2.0 development package

* `apt` users: `libsdl2-dev`
* `yum` users: `SDL2-devel`
* `pacman` users: `sdl2`

If all else fails, the source is available on the [SDL website](https://www.libsdl.org/download-2.0.php#source).

### Install PulseAudio development package

* `apt` users: `libpulse-dev`
* `yum` users: `pulseaudio-libs-devel`
* `pacman` users: `libpulse`

Again, the source is available on the [website](http://www.freedesktop.org/wiki/Software/PulseAudio/Download/).

### Download & install `anna`

Fetch the source into your `$GOPATH`. Note that if `$GOPATH` contains multiple
directories separated by colons, you'll need to replace `$GOPATH` in the
following commands with one of those directories.

    git clone --recursive https://github.com/kierdavis/anna $GOPATH/src/github.com/kierdavis/anna

Then build, making sure to enable the Go 1.5 vendoring experiment flag:

    GO15VENDOREXPERIMENT=1 go install github.com/kierdavis/anna

And run:

    $GOPATH/bin/anna

You may need to check your volume control settings to ensure that `anna` is
capturing audio from the right source (microphone or an output monitor). It's
also a good idea to tweak the recording volume to get the best results.
