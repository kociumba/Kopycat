# Kopycat - Easy to use directory synchronization

Kopycat is a cross-platform directory synchronization tool. It can run as a background service, so you don't have to worry about leaving it open or accidentally closing it. It has all the features you'd expect from a sync tool 😎

## 🚨 Compatibility 🚨

Right now Kopycat has only been tested on Windows, if you want to contribute by testing or improving it on macOS or Linux please [open an issue](https://github.com/kociumba/kopycat/issues).

## Usage

To use Kopycat, download the latest release from [releases](https://github.com/kociumba/kopycat/releases) and run it.

By default, Kopycat will run in terminal mode where it's fully functional and works like any other app.

> [!NOTE]
> Please keep in mind all configuration files and logs are stored relative to the binary, so make sure you aren't running it from `Downloads` 💀

To install Kopycat as a service, pass the `install` flag (if the service is already installed, `install` will try to reinstall it).

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
> 🚨 DO NOT USE KOPYCAT TO SYNC WHOLE DRIVES OR LARGE FOLDERS 🚨
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
- [ ] Control over resource usage
- [ ] Decoupling log cleaning from the start or restart of the app
