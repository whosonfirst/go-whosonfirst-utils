# go-whosonfirst-utils

Go tools and utilities for working with Who's On First documents.

## Tools

### wof-d2fc

Create a GeoJSON `FeatureCollection` from a Git diff (so "diff2featurecollection" or "d2fc"). What that really means is produce FeatureCollection from any list of files formatted per the output of `git diff --name-only`. Basically I wanted a way to easily snapshot one or records in advance of a Git reset so that's what this tool is optimized for.

For example:

```
$> git diff --name-only HEAD..01a6fdd25b7de2d3da7aa2f53f4f44a7efe81c47 | /usr/local/bin/wof-d2fc -repo /usr/local/data/whosonfirst-data-venue-us-ca | python -mjson.tool
{
    "features": [
        {
            "bbox": [
                -122.46228,
                37.783038,
                -122.46228,
                37.783038
            ],

... and so on
```

It is left to users to filter out any non-GeoJSON files from the list passed to `wof-d2fc`.

### wof-expand

_Please write me_

### wof-hash

_Please write me_

## See also

* https://github.com/whosonfirst/go-whosonfirst-uri