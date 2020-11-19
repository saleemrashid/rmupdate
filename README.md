# rmupdate

Utility for fetching software updates from the reMarkable update server.

The reMarkable uses the CoreOS "update_engine" daemon to fetch update payloads from an Omaha server. An update payload contains "installation operations", which describe how to install the update to the device's eMMC.

This utility can fetch and verify these payloads, then reconstruct the installation operations into a root filesystem image.

## Installation

```
go get github.com/saleemrashid/rmupdate/cmd/rmupdate
```

## Usage

The command-line utility has a bunch of configuration options, but typically you want to:

1. Download the latest update payload with `rmupdate fetch -p PLATFORM`, where `PLATFORM` is either `reMarkable` (reMarkable 1) or `reMarkable2` (reMarkable 2).

2. Extract the payload with `rmupdate extract -i PAYLOAD -o ROOTFS` where `PAYLOAD` is the payload filename, and `ROOTFS` is the file you want to write the root filesystem to.

```shell
$ rmupdate fetch -p reMarkable2
2020/12/26 01:35:38 Payload URLs: []string{"https://eu-central-1.linodeobjects.com:443/remarkable-2/build/reMarkable%20Device%20Beta/RM110/2.5.0.27/2.5.0.27_reMarkable2.signed"}
2020/12/26 01:35:38 URL: https://eu-central-1.linodeobjects.com:443/remarkable-2/build/reMarkable%20Device%20Beta/RM110/2.5.0.27/2.5.0.27_reMarkable2.signed
2020/12/26 01:35:38 Size: 66302398 bytes
2020/12/26 01:35:38 Expected SHA-1: 869b8b7e192b4f2a8ccb53395f5f6fcf0e252568
2020/12/26 01:35:38 Expected SHA-256: e32310321499ac70af094c0bf7fedfc80a0cd1e80fe84972fb6de0a6c56f62a7
2020/12/26 01:35:38 Output filename: 2.5.0.27_reMarkable2.signed
Fetching 100% |███████████████████████████████| (63/63 MB, 1.827 MB/s)
2020/12/26 01:36:13 Computed SHA-1: 869b8b7e192b4f2a8ccb53395f5f6fcf0e252568
2020/12/26 01:36:13 Computed SHA-256: e32310321499ac70af094c0bf7fedfc80a0cd1e80fe84972fb6de0a6c56f62a7
$ rmupdate extract -i 2.5.0.27_reMarkable2.signed -o rootfs.ext4
2020/12/26 01:36:39 Parsing manifest
2020/12/26 01:36:39 Verifying RSA signature
2020/12/26 01:36:39 Extracting payload
Extracting 100% |██████████████████████████| (254/254 MB, 22.348 MB/s)
$ file rootfs.ext4
rootfs.ext4: Linux rev 1.0 ext4 filesystem data, UUID=f3fdc696-b818-4348-86ad-67b0812a79fe (extents) (large files) (huge files)
```
