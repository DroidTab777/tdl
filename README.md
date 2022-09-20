## Intro

![](https://img.shields.io/github/go-mod/go-version/iyear/tdl?style=flat-square)
![](https://img.shields.io/github/license/iyear/tdl?style=flat-square)
![](https://img.shields.io/github/workflow/status/iyear/tdl/master%20builder?style=flat-square)
![](https://img.shields.io/github/v/release/iyear/tdl?color=red&style=flat-square)
![](https://img.shields.io/github/last-commit/iyear/tdl?style=flat-square)

📥 Telegram Downloader, but more than a downloader 🚀

> ⚠ Note: Command compatibility is not guaranteed in the early stages of development

> ⚠ Warning: some accounts have been blocked, so please use carefully. Go to [ISSUE](https://github.com/iyear/tdl/issues/21) for discussion

## Features

- Single file start-up
- Low resource usage
- Take up all your bandwidth
- Faster than official clients
- Download files from (protected) chats
- Upload files to Telegram

## Preview

It reaches my proxy's speed limit, and the **speed depends on whether you are a premium**

![](img/preview.gif)

## Install

Go to [GitHub Releases](https://github.com/iyear/tdl/releases) to download the latest version

## Usage

```shell
# get help
tdl -h

# check the version
tdl version

# specify the namespace
tdl -n iyear
# or
export TDL_NS=iyear

# use proxy, only support socks now
tdl --proxy socks5://localhost:1080
# or
export TDL_PROXY=socks5://localhost:1080

# login your account with a name, default is phone & code mode
tdl login -n iyear

# if you have official desktop client on machine, you can import its session
# may be use official session can reduce ban risk(no guarantee)
tdl login -n iyear-desktop -d /path/to/Telegram

# list your chats
tdl chat ls -n iyear

# download files in url mode, url is the message link
tdl dl url -n iyear -u https://t.me/tdl/1 -u https://t.me/tdl/2

# full examples in download url mode
tdl dl url -n iyear --proxy socks5://localhost:1080 -u https://t.me/tdl/1 -u https://t.me/tdl/2 -s 262144 -t 16 -l 3

# upload files to 'Saved Messages', exclude the specified file extensions
tdl up -n iyear -p /path/to/file -p /path -e .so -e .tmp

# full examples in upload mode
tdl up -n iyear --proxy socks5://localhost:1080 -p /path/to/file -p /path -e .so -e .tmp -s 262144 -t 16 -l 3
```

## Env

Avoid typing the same flag values repeatedly every time by setting environment variables.

**Note: The values of all environment variables have a lower priority than flags.**

What flags mean: [flags](docs/command/tdl.md#options)

|    NAME     |      FLAG      |
|:-----------:|:--------------:|
|   TDL_NS    |   `-n/--ns`    |
|  TDL_PROXY  |   `--proxy`    |
|  TDL_DEBUG  |   `--debug`    |
|  TDL_SIZE   |  `-s/--size`   |
| TDL_THREADS | `-t/--threads` |
|  TDL_LIMIT  |  `-l/--limit`  |

## Data

Your account information will be stored in the `~/.tdl` directory.

## Commands

Go to [command documentation](docs/command/tdl.md) for full command docs.

## Contribute

- Better command input
- Better interaction
- Better mode support
- ......

Please provide better suggestions or feedback for the project in the form of [SUBMIT ISSUE](https://github.com/iyear/tdl/issues/new)

## FAQ
**Q: Is this a form of abuse?**

A: No. The download and upload speed is limited by the server side. Since the speed of official clients usually does not reach the account limit, this tool was developed to download files at the highest possible speed.

**Q: Will this result in a ban?**

A: I am not sure. All operations do not involve dangerous actions such as actively sending messages to other people. But it's safer to use an unused account for download and upload operations.

## LICENSE

AGPL-3.0 License
