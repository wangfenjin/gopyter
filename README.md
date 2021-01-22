[![Build Status](https://travis-ci.org/wangfenjin/gopyter.svg?branch=master)](https://travis-ci.org/wangfenjin/gopyter)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/wangfenjin/gopyter/blob/master/LICENSE)

# gopyter - Use Go+ in Jupyter notebooks and nteract

`gopyter` is a Go+ kernel for [Jupyter](http://jupyter.org/) notebooks and [nteract](https://nteract.io/).  It lets you use Go+ interactively in a browser-based notebook or desktop app.  Use `gopyter` to create and share documents that contain live Go+ code, equations, visualizations and explanatory text.  These notebooks, with the live Go+ code, can then be shared with others via email, Dropbox, GitHub and the [Jupyter Notebook Viewer](http://nbviewer.jupyter.org/). Go forth and do data science, or anything else interesting, with Go+ notebooks!

**Acknowledgements** - This project is forked from [gophernotes](https://github.com/wangfenjin/gophernotes) and change the internal interpreter to [Go+](https://github.com/goplus/gop).

- [Examples](#examples)
- Install gopyter:
  - [Prerequisites](#prerequisites)
  - [FreeBSD](#linux-or-freebsd)
  - [Linux](#linux-or-freebsd)
  - [Mac](#mac)
  - [Windows](#windows)
  - [Docker](#docker)
- [Getting Started](#getting-started)
- [Limitations](#limitations)
- [Troubleshooting](#troubleshooting)

## Examples

### Jupyter Notebook:

![](files/jupyter.gif)

### nteract:

![](files/nteract.gif)

### Example Notebooks (download and run them locally, follow the links to view in Github, or use the [Jupyter Notebook Viewer](http://nbviewer.jupyter.org/)):
- [Tutorial](examples/Tutorial.ipynb)

## Installation

### Prerequisites

- [Go 1.13+](https://golang.org/doc/install) - including GOPATH/bin added to your PATH (i.e., you can run Go binaries that you `go install`).
- [Jupyter Notebook](http://jupyter.readthedocs.io/en/latest/install.html) or [nteract](https://nteract.io/desktop)
- [git](https://git-scm.com/download) - usually already present on Linux and Mac OS X. If not present, follow the instructions at [https://git-scm.com/download](https://git-scm.com/download)

### Linux or FreeBSD

The instructions below should work both on Linux and on FreeBSD.

Method 1: quick installation as module
```sh
$ env GO111MODULE=on go get github.com/wangfenjin/gopyter
$ mkdir -p ~/.local/share/jupyter/kernels/gopyter
$ cd ~/.local/share/jupyter/kernels/gopyter
$ cp "$(go env GOPATH)"/pkg/mod/github.com/wangfenjin/gopyter@v0.7.1/kernel/*  "."
$ chmod +w ./kernel.json # in case copied kernel.json has no write permission
$ sed "s|gopyter|$(go env GOPATH)/bin/gopyter|" < kernel.json.in > kernel.json
```

Method 2: manual installation from GOPATH
```sh
$ env GO111MODULE=off go get -d -u github.com/wangfenjin/gopyter
$ cd "$(go env GOPATH)"/src/github.com/wangfenjin/gopyter
$ env GO111MODULE=on go install
$ mkdir -p ~/.local/share/jupyter/kernels/gopyter
$ cp kernel/* ~/.local/share/jupyter/kernels/gopyter
$ cd ~/.local/share/jupyter/kernels/gopyter
$ chmod +w ./kernel.json # in case copied kernel.json has no write permission
$ sed "s|gopyter|$(go env GOPATH)/bin/gopyter|" < kernel.json.in > kernel.json
```

To confirm that the `gopyter` binary is installed in GOPATH, execute it directly:
```sh
$ "$(go env GOPATH)"/bin/gopyter
```
and you shoud see the following:
```sh
2017/09/20 10:33:12 Need a command line argument specifying the connection file.
```

**Note** - if you have the `JUPYTER_PATH` environmental variable set or if you are using an older version of Jupyter, you may need to copy this kernel config to another directory.  You can check which directories will be searched by executing:

```sh
$ jupyter --data-dir
```

### Mac

Method 1: quick installation as module
```sh
$ env GO111MODULE=on go get github.com/wangfenjin/gopyter
$ mkdir -p ~/Library/Jupyter/kernels/gopyter
$ cd ~/Library/Jupyter/kernels/gopyter
$ cp "$(go env GOPATH)"/pkg/mod/github.com/wangfenjin/gopyter@v0.7.1/kernel/*  "."
$ chmod +w ./kernel.json # in case copied kernel.json has no write permission
$ sed "s|gopyter|$(go env GOPATH)/bin/gopyter|" < kernel.json.in > kernel.json
```

Method 2: manual installation from GOPATH
```sh
$ env GO111MODULE=off go get -d -u github.com/wangfenjin/gopyter
$ cd "$(go env GOPATH)"/src/github.com/wangfenjin/gopyter
$ env GO111MODULE=on go install
$ mkdir -p ~/Library/Jupyter/kernels/gopyter
$ cp kernel/* ~/Library/Jupyter/kernels/gopyter
$ cd ~/Library/Jupyter/kernels/gopyter
$ chmod +w ./kernel.json # in case copied kernel.json has no write permission
$ sed "s|gopyter|$(go env GOPATH)/bin/gopyter|" < kernel.json.in > kernel.json
```

To confirm that the `gopyter` binary is installed in GOPATH, execute it directly:
```sh
$ "$(go env GOPATH)"/bin/gopyter
```
and you shoud see the following:
```sh
2017/09/20 10:33:12 Need a command line argument specifying the connection file.
```

**Note** - if you have the `JUPYTER_PATH` environmental variable set or if you are using an older version of Jupyter, you may need to copy this kernel config to another directory.  You can check which directories will be searched by executing:

```sh
$ jupyter --data-dir
```

### Windows

1. Copy the kernel config:

    ```
    mkdir %APPDATA%\jupyter\kernels\gopyter
    xcopy %GOPATH%\src\github.com\wangfenjin\gopyter\kernel %APPDATA%\jupyter\kernels\gopyter /s
    ```

    Note, if you have the `JUPYTER_PATH` environmental variable set or if you are using an older version of Jupyter, you may need to copy this kernel config to another directory.  You can check which directories will be searched by executing:

    ```
    jupyter --data-dir
    ```

2. Update `%APPDATA%\jupyter\kernels\gopyter\kernel.json` with the FULL PATH to your gopyter.exe (in %GOPATH%\bin), unless it's already on the PATH.  For example:

    ```
    {
        "argv": [
          "C:\\gopath\\bin\\gopyter.exe",
          "{connection_file}"
          ],
        "display_name": "Go",
        "language": "go",
        "name": "go"
    }
    ```

### Docker

You can try out or run Jupyter + gopyter without installing anything using Docker. To run a Go notebook that only needs things from the standard library, run:

```
$ docker run -it -p 8888:8888 wangfenjin/gopyter
```

Or to run a Go notebook with access to common Go data science packages (gonum, gota, golearn, etc.), run:

```
$ docker run -it -p 8888:8888 wangfenjin/gopyter:latest-ds
```

In either case, running this command should output a link that you can follow to access Jupyter in a browser. Also, to save notebooks to and/or load notebooks from a location outside of the Docker image, you should utilize a volume mount.  For example:

```
$ docker run -it -p 8888:8888 -v /path/to/local/notebooks:/path/to/notebooks/in/docker wangfenjin/gopyter
```

## Getting Started

### Jupyter

- If you completed one of the local installs above (i.e., not the Docker install), start the jupyter notebook server:

  ```
  jupyter notebook
  ```

- Select `Go+` from the `New` drop down menu.

- Have fun!

### nteract

- Launch nteract.

- From the nteract menu select Language -> Go.

- Have fun!

## Limitations

gopyter uses [gop](https://github.com/goplus/gop) under the hood to evaluate Go code interactively. It can only support the code same as GoPlus.  Most notably, gopyter does NOT support:

- import multiple times
- import external packages. You need to follow this [wiki](https://github.com/goplus/gop/wiki/Import-Go-packages-in-GoPlus-programs) page to use other github packages.

## Troubleshooting

### gopyter not found

Depending on your environment, you may need to manually change the path to the `gopyter` executable in `kernel/kernel.json` before copying it to `~/.local/share/jupyter/kernels/gopyter`.  You can put the **full path** to the `gopyter` executable here, and you shouldn't have any further issues.

### "Kernel error" in a running notebook

```
Traceback (most recent call last):
  File "/usr/local/lib/python2.7/site-packages/notebook/base/handlers.py", line 458, in wrapper
    result = yield gen.maybe_future(method(self, *args, **kwargs))
  File "/usr/local/lib/python2.7/site-packages/tornado/gen.py", line 1008, in run
    value = future.result()
  ...
  File "/usr/local/Cellar/python/2.7.11/Frameworks/Python.framework/Versions/2.7/lib/python2.7/subprocess.py", line 1335, in _execute_child
    raise child_exception
OSError: [Errno 2] No such file or directory
```

Stop jupyter, if it's already running.

Add a symlink to `/go/bin/gopyter` from your path to the gopyter executable. If you followed the instructions above, this will be:

```
sudo ln -s $HOME/go/bin/gopyter /go/bin/gopyter
```

Restart jupyter, and you should now be up and running.

### error "could not import C (no metadata for C)" when importing a package

At a first analysis, it seems to be a limitation of the new import mechanism that supports Go modules.
You can switch the old (non module-aware) mechanism with the command `%go111module off`

To re-enable modules support, execute `%go111module on`

### Look at Jupyter notebook's logs for debugging

In order to see the logs for your Jupyter notebook, use the --log-level option
```
jupyter notebook --log-level DEBUG
```
