# iputil

Utility functions for IP address for Go

## Installation

```
$ go get github.com/mattn/iputil
```

## Tips

Ban IP range from an IP address.

```
$ iprange XXX.XXX.XXX.XXX | xargs -n 1 sudo ufw insert 1 deny from
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
