NAME := alp
GO_LDFLAGS=-s -w
GIT_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build -trimpath -ldflags "$(GO_LDFLAGS) -X=main.version=$(GIT_COMMIT)" -o $(NAME) ./cmd/$(NAME)

run_nginx:
	docker run -p 18080:80 -v $(PWD)/tmp:/var/log/nginx -v $(PWD)/dockerfiles/nginx/nginx.conf:/etc/nginx/nginx.conf --rm -t nginx:latest

run_h2o_ltsv:
	docker run -p 18080:8080 -v $(PWD)/tmp:/home/h2o/logs -v $(PWD)/dockerfiles/h2o/ltsv.h2o.conf:/home/h2o/h2o.conf --rm -t lkwg82/h2o-http2-server

run_h2o_json:
	docker run -p 18080:8080 -v $(PWD)/tmp:/home/h2o/logs -v $(PWD)/dockerfiles/h2o/json.h2o.conf:/home/h2o/h2o.conf --rm -t lkwg82/h2o-http2-server

run_h2o_regexp:
	docker run -p 18080:8080 -v $(PWD)/tmp:/home/h2o/logs -v $(PWD)/dockerfiles/h2o/regexp.h2o.conf:/home/h2o/h2o.conf --rm -t lkwg82/h2o-http2-server

run_httpd_ltsv:
	docker run -p 18080:80 -v $(PWD)/tmp:/usr/local/apache2/logs -v $(PWD)/dockerfiles/apache/ltsv.httpd.conf:/usr/local/apache2/conf/httpd.conf --rm -t httpd:latest

run_httpd_json:
	docker run -p 18080:80 -v $(PWD)/tmp:/usr/local/apache2/logs -v $(PWD)/dockerfiles/apache/json.httpd.conf:/usr/local/apache2/conf/httpd.conf --rm -t httpd:latest

run_httpd_regexp:
	docker run -p 18080:80 -v $(PWD)/tmp:/usr/local/apache2/logs -v $(PWD)/dockerfiles/apache/regexp.httpd.conf:/usr/local/apache2/conf/httpd.conf --rm -t httpd:latest
