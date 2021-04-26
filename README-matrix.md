# [PROG_NAME] vXX.X

#### Last update 202x-xx-xx

---

##### This program is designed for . . .

###### Why ?

> . . .

---

###### Requirements:

- **libgtksourceview-4-0** is required to work (most Linux distributions include it natively)

```bash
$ sudo apt install libgtksourceview-4-0
```

---

- If you just want **to use it**, simply **download** the **compiled version** ( .deb) under the [Releases](https://github.com/hfmrow/%5BPROG_NAME%5D/releases) tab.

- Otherwise, if you plan to play inside the source code, see below **How to compile** section.

- All suggestions, contributions and ideas to improve software usability will be greatly appreciated.

**[PROG_NAME]** Debian package installation:

> ```bash
> $ sudo dpkg -i [PROG_NAME]-x.x-amd64.deb
> ```
> 
> Uninstall:
> 
> ```bash
> $ sudo dpkg -P [PROG_NAME]
> ```

---

### How it's made

- Programmed with go language: [golang](https://golang.org/doc/)
- GUI provided by [Gotk3](https://github.com/gotk3/gotk3), GUI library for Go (minimum required gtk3.16).
- Text editor use [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) component.
- Version check (for update checks), [go-curl](https://github.com/andelf/go-curl), to connect and retrieve file from endpoint **https**. In my opinion, this is a better way (less memory consuming) rather than using the native Go **http** package for this kind of operation...
  **Note**: only *get* https method is used (to get file that contain latest version information).
- I use home-made software: "Gotk3ObjHandler" to embed images/icons, UI-information and manage/generate gtk3 objects code from [glade ui designer](https://glade.gnome.org/). and "Gotk3ObjTranslate" to generate the language files and the assignment of a tool-tip on the gtk3 objects (both are not published at the moment, in fact, they need documentations and, I have not had the time to do them).

### Functionalities

- 1
- 2
- 3
- ...

### Some pictures

![whole.jpg](assets/readme/)

### How to compile

1. Be sure you have [Golang version >= 1.16](https://golang.org/dl/) installed in right way, [Go installation](https://golang.org/doc/install).

2. Getting source from repository:
   
   1. **git preferred method (the most close to $GOPATH old usage):**      
   
   ```bash
   $ cd "your-local-workspace"
   $ git clone "https://github.com/hfmrow/[PROG_NAME]" "[PROG_NAME]"
   $ cd [PROG_NAME]
   $ go build . && ./[PROG_NAME]
   ```
   
   2. **golang method (this one put retrieved package to  $GOPATH/pkg/pkgName):**    
   
   ```bash
   $ go get -d "github.com/hfmrow/[PROG_NAME]"
   $ cd ${GOPATH}/pkg/mod/github.com/hfmrow/[PROG_NAME]*
   $ go build . && ./[PROG_NAME]
   ```

3. Ensure that you Have [libcurl](https://curl.se/libcurl/) (`C` headers) correctly installed.
   
   ```bash
   $ # For Ubuntu:
   $ sudo apt install libcurl4-gnutls-dev
   ```

##### If you have gomodules enabled, all of the following procedure will be done automatically you can skip these.

- Install https://github.com/andelf/go-curl. Assuming you have already done the previous 'step 3'.

- Install [Go bindings for GTK3](https://github.com/gotk3/gotk3) and follow [Installation instructions](https://github.com/gotk3/gotk3/wiki#installation).

- Install [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) library and follow instructions to install it with his `libgtksourceview-X-dev` (`C`headers) package.

---

### Playing with code

- Open terminal window and at command prompt, type: `go get github.com/hfmrow/[PROG_NAME]`

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

<details>
  <summary>Built using</summary>

| Name                                                       | Version / Info / Name                       |
| ---------------------------------------------------------- | ------------------------------------------- |
| GOLANG                                                     | V1.16.3 -> GO111MODULE="on", GOPROXY="auto" |
| DISTRIB                                                    | LinuxMint Xfce                              |
| VERSION                                                    | 20.1                                        |
| CODENAME                                                   | ulyana                                      |
| RELEASE                                                    | #46-Ubuntu SMP Fri Jul 10 00:24:02 UTC 2020 |
| UBUNTU_CODENAME                                            | focal                                       |
| KERNEL                                                     | 5.8.0-50-generic                            |
| HDWPLATFORM                                                | x86_64                                      |
| GTK+ 3                                                     | 3.24.20                                     |
| GLIB 2                                                     | 2.64.3                                      |
| CAIRO                                                      | 1.16.0                                      |
| [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) | 4.6.0                                       |
| [LiteIDE](https://github.com/visualfc/liteide)             | 37.4 qt5.x                                  |
| Qt5                                                        | 5.12.8 in /usr/lib/x86_64-linux-gnu         |

</details>

- The compilation have not been tested under Windows or Mac OS.

### You got an issue ?

- Go to this page: [Issues](https://github.com/hfmrow/[PROG_NAME]/issues) and start a new problem report.
- Give the information (as above), concerning your working environment as well as the version of the operating system used.
- Provide a method to reproduce the problem.

### Used libraries

- [Go bindings for GTK3](https://github.com/gotk3/gotk3)
- [GitHub - andelf/go-curl: golang curl(libcurl) binding](https://github.com/andelf/go-curl)
- [Golang GtkSourceView binding for use with gotk3](https://github.com/hfmrow/gotk3_gtksource)
- And some personal libraries not yet published.

### Visit

- [GitHub's hfmrow (H.F.M) Repositories](https://github.com/hfmrow?tab=repositories).
- [Website hfmrow's Linux softwares](https://hfmrow.go.yo.fr/)
