# URLsplit

A tool to enable script processing of URLs


### Install

#### Homebrew

```
brew tap relvacode/urlsplit
brew install urlsplit
```

#### Linux / Docker

```
curl --fail -L https://github.com/relvacode/urlsplit/releases/download/0.2/urlsplit_0.2_linux_x86_64.tar.gz | sudo tar -xzf - -C /usr/bin
```

### Usage

#### Export Environment Variables

You can export a set of environment variables using the `-e` option.
This will print a set of export statements that can be evaluated by your shell.

```
urlsplit -e <URL>
```

```
$ urlsplit -e 'http://username:password@example.org/path?query=yes'
export "URL_SCHEME=http"
export "URL_HOST=example.org"
export "URL_HOSTNAME=example.org"
export "URL_PORT="
export "URL_USERNAME=username"
export "URL_PASSWORD=password"
export "URL_URI=/path?query=yes"
export "URL_PATH=/path"
export "URL_ESCAPED_PATH=/path"
export "URL_QUERY=query=yes"
export "URL_FRAGMENT="
export "URL_QUERY_query=yes"
```

#### Get a Specific Key

Use the `-k` option combined with a key name to print that specific value.

```
urlsplit -k <KEY> <URL>
```

```
$ urlsplit -k URL_HOST 'http://username:password@example.org/path?query=yes'
example.org
```

#### Render a Template

Using the `-f` option you can render a Django-style template from URL parameters.
Useful for printing conditionally or generating command line arguments.

Check out [the Django template language](https://docs.djangoproject.com/en/dev/topics/templates/#the-django-template-language) for details on how to use this syntax.

```
urlsplit -f <TEMPLATE> <URL>
```

```
$ urlsplit -f 'pg_isready -h {{URL_HOSTNAME}}{% if URL_PORT %} -p {{URL_PORT}}{% endif %}{% if URL_USERNAME %} -U {{URL_USERNAME}}{% endif %}' 'postgres://postgres:root@localhost/postgres'
pg_isready -h localhost -U postgres
```