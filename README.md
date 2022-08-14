
# cms-lookup

Find the CMS used by a given site.

### Build
- Install [go](https://go.dev/)
- Run `go build .` inside the project directory


### Usage
```
./cms_lookup:
  -filename string
        File containing the list of URL's to process.
  -threads int
        Number of threads to run. (default 5)
```