# Go zipext

![CI Linux](https://github.com/enr/zipext/workflows/CI%20Nix/badge.svg)
![CI Windows](https://github.com/enr/zipext/workflows/CI%20Windows/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/enr/zipext)](https://pkg.go.dev/github.com/enr/zipext)
[![Go Report Card](https://goreportcard.com/badge/github.com/enr/zipext)](https://goreportcard.com/report/github.com/enr/zipext)


Go library to manipulate zip files.

Import library:

```Go
    import (
        "github.com/enr/zipext"
    )
```

Create a zip archive:

```Go
    contents := "/path/to/files"
    zipPath := "/path/to/archive.zip"
    err := zipext.Create(contents, zipPath)
    if err != nil {
        t.Errorf("error in Create(%s,%s): %s %s", contents, zipPath, reflect.TypeOf(err), err.Error())
    }
```

Extract contents from zip:

```Go
    err = zipext.Extract(zipPath, extractPath)
    if err != nil {
        t.Errorf("error in Extract(%s,%s): %s %s", zipPath, unzipDir, reflect.TypeOf(err), err.Error())
    }
```

Visit zip contents:

```Go
    err := zipext.Walk(path, func(f *zip.File, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("- %#v\n", f)
        return nil
    })
```

Check if file is valid zip:

```Go
    p := "/path/to/archive.zip"
    valid, err := zipext.IsValidZip(p)
    if err != nil {
        t.Errorf("got error checking for valid zip '%s'", p)
    }
```

## License

Apache 2.0 - see LICENSE file.

Copyright 2014-TODAY zipext contributors
