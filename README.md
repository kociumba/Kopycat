# Kopycat - Easy to use directory synchronization

<p align="center">
  <img src="/winres/icon.png" alt="Icon">
</p>


[![Codacy Badge](https://app.codacy.com/project/badge/Grade/a852d0bff81d476caca531dc2ebb3a14)](https://app.codacy.com/gh/kociumba/Kopycat/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/kociumba/kopycat/latest)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/kociumba/kopycat/total)
[![Release](https://github.com/kociumba/Kopycat/actions/workflows/release.yml/badge.svg)](https://github.com/kociumba/Kopycat/actions/workflows/release.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kociumba/kopycat)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/kociumba)

Kopycat is a cross-platform directory synchronization tool. It can run as a background service, so you don't have to worry about leaving it open or accidentally closing it. It has all the features you'd expect from a sync tool 😎

## 🚨 Compatibility and usability 🚨

Kopycat is functional in its current state but is not yet production-ready. When it meets the planned ease of use and readiness, it will be released as v1.0.0.

If the current release version is lower than v1.0.0, be prepared for potential breaking changes and difficulties when using Kopycat.

**Right now Kopycat has only been tested on Windows, if you want to contribute by testing or improving it on macOS or Linux please [open an issue](https://github.com/kociumba/kopycat/issues).**

## Usage

To use Kopycat, download the latest release from the [releases page](https://github.com/kociumba/kopycat/releases) and run it.

By default, Kopycat runs in terminal mode where it is fully functional and operates like any other application.

> [!NOTE]
> Please keep in mind all configuration files and logs are stored relative to the binary, so make sure you aren't running it from `Downloads` 💀

To install Kopycat as a service, use the `install` flag (if the service is already installed, `install` will attempt to reinstall it).

> [!WARNING]
> This will require admin privileges.
> On Windows, run in an elevated shell, and on Linux, use `sudo`.

From here, you can pass these flags to Kopycat with admin privileges to manage the installed service:

- `start` to start the service (`install` automatically starts the service on install)
- `restart` to restart the service if there are any issues
- `stop` to stop the service without restarting it
- `remove` to stop and uninstall the service

After running Kopycat in any mode, you can open the locally hosted web GUI at http://localhost:42069 to configure it.

The port is set to `42069` by default but can be changed with the `-port` flag.

> [!IMPORTANT]
> 🚨 DO NOT USE KOPYCAT TO SYNC WHOLE DRIVES OR LARGE DIRECTORIES 🚨
> 
> It was not designed to handle these cases and may malfunction or consume a lot of system resources.

### What to sync ⚙️

| Example sync targets                              | For use with Kopycat? |
| :------------------------------------------------ | :-------------------: |
| ~/ *(Obviously)*   | ❌ |
| ~/Pictures *(Typically large)*  | ❌ |
| ~/Videos *(Typically large)*    | ❌ |
| ~/.config *(Small things like configs)* | ✅ |

## Note about releases

**I do not own a Mac nor do I have Xcode. Because of this, [releases](https://github.com/kociumba/kopycat/releases) contain a Mach-O binary for macOS but it is not packaged as a `.app` or `.dmg`.**

## Project Scope

This isn't a complete list by any means, but it should serve as a baseline of what you can expect from Kopycat.

- [x] Run as a service (tested on Windows, theoretically should be compatible with Linux and macOS)
- [x] Locally hosted web GUI for configuration
- [x] Terminal mode when running as a normal app
- [ ] Control over resource usage (this is a hard one and I'm not entirely sure I will be able to achive it in GO)
- [x] Decoupling log cleaning from the start or restart of the app

## Final notes

At the moment, the web GUI won't function without an internet connection as it uses [HTMX](https://htmx.org/) 
and [jQuery](https://jquery.com/).

I used this to make development much faster and easier.

The rest of the app remains functional.
The only caveat being you would have to configure it by sending manual server requests or edit the `config/config.json` file.
