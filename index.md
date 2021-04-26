# 7z-folder v1.6

#### Last update 2021-04-26

---

##### This program is designed for . . .

Compressing folder and files to 7z format. Some options are available. The 7za used command store filesystem permissions (such as UNIX owner/group permissions or NTFS ACLs). Not designed for large backup/archival purposes. On Ubuntu, use 'sudo apt-get install p7zip-full' to install required command.

###### What for:

Quick backup of the local repository in case of problem with the new modifications (for example).

---

###### Requirements:

- **p7zip-full** is required to work.

```bash
$ # On Ubuntu, use: 
$ sudo apt-get install p7zip-full
```

---

- If you just want **to use it**, simply **download** the **compiled version** ( .deb) under the [Releases](https://github.com/hfmrow/7z-folder/releases) tab.

- Otherwise, if you plan to play inside the source code, see below **How to compile** section.

- All suggestions, contributions and ideas to improve software usability will be greatly appreciated.

**7z-folder** Debian package installation:

> ```bash
> $ sudo dpkg -i 7z-folder-1.6-amd64.deb
> ```
> 
> Uninstall:
> 
> ```bash
> $ sudo dpkg -P 7z-folder
> ```

---

### How it's made

- Programmed with go language: [golang](https://golang.org/doc/)
- GUI provided by [Gotk3](https://github.com/gotk3/gotk3), GUI library for Go (minimum required gtk3.16).
- I use home-made software: "Gotk3ObjHandler" to embed images/icons, UI-information and manage/generate gtk3 objects code from [glade ui designer](https://glade.gnome.org/). and "Gotk3ObjTranslate" to generate the language files and the assignment of a tool-tip on the gtk3 objects (both are not published at the moment, in fact, they need documentations and, I have not had the time to do them).

### Functionalities

- Auto increment filename [00], [01], [02]...
- Parametrable compression lvl and dictionary size.
- Append, Upate, New archive options.
- D&D capable.

### Some pictures

![whole.jpg](assets/readme/main.jpg)

### How to compile

1. Be sure you have [Golang version >= 1.16](https://golang.org/dl/) installed in right way, [Go installation](https://golang.org/doc/install).

2. Getting source from repository:
   
   1. **git preferred method (the most close to $GOPATH old usage):**      
   
   ```bash
   $ cd "your-local-workspace"
   $ git clone "https://github.com/hfmrow/7z-folder" "7z-folder"
   $ cd 7z-folder
   $ go build . && ./7z-folder
   ```
   
   2. **golang method (this one put retrieved package to  $GOPATH/pkg/pkgName):**    
   
   ```bash
   $ go get -d "github.com/hfmrow/7z-folder"
   $ cd ${GOPATH}/pkg/mod/github.com/hfmrow/7z-folder*
   $ go build . && ./7z-folder
   ```

##### If you have gomodules enabled, all of the following step will be done automatically you can skip it.

- Install [Go bindings for GTK3](https://github.com/gotk3/gotk3) and follow [Installation instructions](https://github.com/gotk3/gotk3/wiki#installation).

---

### Playing with code

- **Since** [Golang version >= 1.16](https://golang.org/dl/), native `embed` library is used instead of previous one and the following restriction (striked) is out of date. ~~To change gtk3 interface you need to set `devMode` flag at `true`.~~ A home made software, (not published actually) have been used to generate some parts of source code / assets embedding. ~~So, you cannot (at this time) change interface for production mode.~~

- To change language file quickly, you need to use another home made software, (not published actually). You can still do it manually, all data has been stored in a '.json' file in the `assets/lang` directory.

- To Produce a *stand-alone executable*, you must change inside `main.go` file:

```go
    func main() {
        devMode = true
    ...    
```

into

```go
    func main() {
        devMode = false
    ...
```

This operation indicates that the internal behavior of the software will be modified to adapt to the production environment (display of errors, location of the configuration file, etc.).

## Os information:

| Name                                                       | Version / Info / Name                               |
| ---------------------------------------------------------- | --------------------------------------------------- |
| GOLANG                                                     | V1.16.3 -> GO111MODULE="on", GOPROXY="auto"         |
| DISTRIB                                                    | LinuxMint Xfce                                      |
| VERSION                                                    | 20.1                                                |
| CODENAME                                                   | ulyssa                                              |
| RELEASE                                                    | #56~20.04.1-Ubuntu SMP Mon Apr 12 21:46:35 UTC 2021 |
| UBUNTU_CODENAME                                            | focal                                               |
| KERNEL                                                     | 5.8.0-50-generic                                    |
| HDWPLATFORM                                                | x86_64                                              |
| GTK+ 3                                                     | 3.24.20                                             |
| GLIB 2                                                     | 2.64.6                                              |
| CAIRO                                                      | 1.16.0                                              |
| [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) | 4.6.0                                               |
| [LiteIDE](https://github.com/visualfc/liteide)             | 37.4 qt5.x                                          |
| Qt5                                                        | 5.12.8 in /usr/lib/x86_64-linux-gnu                 |

- The compilation have not been tested under Windows or Mac OS.

### You got an issue ?

- Go to this page: [Issues](https://github.com/hfmrow/7z-folder/issues) and start a new problem report.
- Give the information (as above), concerning your working environment as well as the version of the operating system used.
- Provide a method to reproduce the problem.

### Used libraries

- [Go bindings for GTK3](https://github.com/gotk3/gotk3)
- And some personal libraries not yet published.

### Visit

- [GitHub's hfmrow (H.F.M) Repositories](https://github.com/hfmrow?tab=repositories).
- [Website hfmrow's Linux softwares](https://hfmrow.go.yo.fr/)
