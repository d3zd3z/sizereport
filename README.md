# Size report

This is a simple size comparator for elf files.  Currently, it is hard
coded for the tools and a few other things.

## Getting it

Once $GOPATH is setup, run:

```bash
$ go get github.com/d3zd3z/sizereport
```

which should place a sizereport executable in $GOPATH/bin

## Running it

You can give it a single elf file, which will print out the sizes of
each symbol, or two elf files, which will print out what changed (size
wise) between the two elf files.
