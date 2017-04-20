# go-whosonfirst-utils

Go tools and utilities for working with Who's On First documents.

## Tools

### wof-compare

```
./bin/wof-compare -h
Usage of ./bin/wof-compare:
  -filelist
    	Read WOF IDs from a "file list" document.
```

Compare one or more Who's On First documents against multiple sources, reading IDs from the command line or a "file list" document. Current sources are:

* wof _https://whosonfirst.mapzen.com/data/_
* github _https://raw.githubusercontent.com/whosonfirst-data/whosonfirst-data/master/data/_
* s3 _https://s3.amazonaws.com/whosonfirst.mapzen.com/data/_

For example:

```
./bin/wof-compare 890537261 85898837
wofid,match,github,wof,s3
890537261,MATCH,2e4d929a9faad66043f3c5790363ba3e,2e4d929a9faad66043f3c5790363ba3e,2e4d929a9faad66043f3c5790363ba3e
85898837,MATCH,b64b923c27d68abfddd234ee6c8c584a,b64b923c27d68abfddd234ee6c8c584a,b64b923c27d68abfddd234ee6c8c584a
```

You can also pass a "file list" document which is pretty much what it sounds like: A list of WOF documents, one per line. For example here is how you might generate a file list for all the descendants of Helsinki using the [wof-api](https://github.com/whosonfirst/go-whosonfirst-api#wof-api) tool:

```
./bin/wof-api -param method=whosonfirst.places.getDescendants -param id=101748417 -param api_key=mapzen-xxxxxx -filelist -filelist-prefix /usr/local/data/whosonfirst-data/data -paginated 
/usr/local/data/whosonfirst-data/data/858/988/21/85898821.geojson
/usr/local/data/whosonfirst-data/data/858/988/25/85898825.geojson
/usr/local/data/whosonfirst-data/data/858/988/37/85898837.geojson
/usr/local/data/whosonfirst-data/data/858/988/43/85898843.geojson
/usr/local/data/whosonfirst-data/data/858/988/45/85898845.geojson
/usr/local/data/whosonfirst-data/data/858/988/47/85898847.geojson
/usr/local/data/whosonfirst-data/data/858/988/51/85898851.geojson
/usr/local/data/whosonfirst-data/data/858/988/55/85898855.geojson
/usr/local/data/whosonfirst-data/data/858/988/61/85898861.geojson
/usr/local/data/whosonfirst-data/data/858/988/65/85898865.geojson
/usr/local/data/whosonfirst-data/data/858/988/67/85898867.geojson
/usr/local/data/whosonfirst-data/data/858/988/75/85898875.geojson
/usr/local/data/whosonfirst-data/data/858/988/81/85898881.geojson
/usr/local/data/whosonfirst-data/data/858/988/85/85898885.geojson
... and so on
```

Let's imagine you wrote those results to a file called `helsinki.txt`. You can compare all the WOF IDs listed like this:

```
./bin/wof-compare -filelist helsinki.txt 
wofid,match,github,s3,wof
85898821,MATCH,9b0b76f328a64d9f0eb3ab604ea855eb,9b0b76f328a64d9f0eb3ab604ea855eb,9b0b76f328a64d9f0eb3ab604ea855eb
85898825,MATCH,a799cc0f2d46c3b0c94c2f014286b16d,a799cc0f2d46c3b0c94c2f014286b16d,a799cc0f2d46c3b0c94c2f014286b16d
85898837,MATCH,b64b923c27d68abfddd234ee6c8c584a,b64b923c27d68abfddd234ee6c8c584a,b64b923c27d68abfddd234ee6c8c584a
85898843,MATCH,c557962d544f66c51777629568dcb542,c557962d544f66c51777629568dcb542,c557962d544f66c51777629568dcb542
85898845,MATCH,6e7e8e5ab00b8e92bc3f9c1de512c866,6e7e8e5ab00b8e92bc3f9c1de512c866,6e7e8e5ab00b8e92bc3f9c1de512c866
85898847,MATCH,0822020aaed4b53379df793cbaebe2b5,0822020aaed4b53379df793cbaebe2b5,0822020aaed4b53379df793cbaebe2b5
85898851,MATCH,9dc34d1ef12f38e329b5e8b3f8bd8533,9dc34d1ef12f38e329b5e8b3f8bd8533,9dc34d1ef12f38e329b5e8b3f8bd8533
85898855,MATCH,8f099ef74842b1ec42990d7a4dc165a8,8f099ef74842b1ec42990d7a4dc165a8,8f099ef74842b1ec42990d7a4dc165a8
85898861,MATCH,38f0fb958a9fd997b22046923e9fa443,38f0fb958a9fd997b22046923e9fa443,38f0fb958a9fd997b22046923e9fa443
85898865,MATCH,989d91f616a4a71fcc6c4b4da0850d4c,989d91f616a4a71fcc6c4b4da0850d4c,989d91f616a4a71fcc6c4b4da0850d4c
85898867,MATCH,94ad036fcf2bd7b1c0f8bda661e50e2f,94ad036fcf2bd7b1c0f8bda661e50e2f,94ad036fcf2bd7b1c0f8bda661e50e2f
85898875,MATCH,ec3114af9ef10bf18c9e01a016db31d8,ec3114af9ef10bf18c9e01a016db31d8,ec3114af9ef10bf18c9e01a016db31d8
85898881,MATCH,c6f576fee8e47d06eb6b9baae2dffe89,c6f576fee8e47d06eb6b9baae2dffe89,c6f576fee8e47d06eb6b9baae2dffe89
85898885,MATCH,5c16f782f59366e52b8eb3f07dbdac10,5c16f782f59366e52b8eb3f07dbdac10,5c16f782f59366e52b8eb3f07dbdac10
85898887,MATCH,80a8bb00344e65f69bc5e0891e45f491,80a8bb00344e65f69bc5e0891e45f491,80a8bb00344e65f69bc5e0891e45f491
85898895,MATCH,7f4d38818ed11cb5f484557a6e8b723c,7f4d38818ed11cb5f484557a6e8b723c,7f4d38818ed11cb5f484557a6e8b723c
85898897,MATCH,706b9266f7815f79834b410189a7e917,706b9266f7815f79834b410189a7e917,706b9266f7815f79834b410189a7e917
... and so on
```

#### Caveats

* See the way one of the sources is `https://raw.githubusercontent.com/whosonfirst-data/whosonfirst-data/master/data/`. That means if a WOF ID is part of a _different_ repo (say `whosonfirst-data-venue-fi`) it won't be able to be found so the comparison will always fail. We should figure out how to fix that, but today we can't. One way would to be allow query results against the [Who's On First API](https://mapzen.com/documentation/wof) to be filtered by `wof:repo` but again that's not possible today... 

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