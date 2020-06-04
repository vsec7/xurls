# xurls

A URLs eXtractor from source.

## Install
```
 ▶ go get -u github.com/vsec7/xurls
```

## Usage

```
ve@cans ~$ xurls --help 

xURLs (eXtract URLs)

By : viloid [Sec7or - Surabaya Hacker Link]

Basic Usage :
 ▶ echo http://domain.com/path/file.js | xurls
 ▶ cat listurls.txt | xurls -o result.txt

Options :
  -H, --header <header>                 Header to the request
  -o, --output <output>                 Output file (*default xurls.txt)
  -x, --proxy <proxy>   				HTTP proxy
```

## Credit And Thanks
```
@tomnomnom (github.com/tomnomnom)
```