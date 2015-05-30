Go zipext
=========

[![Build Status](https://travis-ci.org/enr/go-zipext.png?branch=master)](https://travis-ci.org/enr/go-zipext)
[![Build status](https://ci.appveyor.com/api/projects/status/2nnl8sqg31b9vrvm?svg=true)](https://ci.appveyor.com/project/enr/go-zipext)

Go library to manipulate zip files.

Import library:

```Go
    import (
        "github.com/enr/go-zipext/zipext"
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


License
-------

Apache 2.0 - see LICENSE file.

   Copyright 2014 go-zipext contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
