prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-utils; then rm -rf src/github.com/whosonfirst/go-whosonfirst-utils; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-utils
	cp utils.go src/github.com/whosonfirst/go-whosonfirst-utils/utils.go
