# Paste Table

(Work-in-progress)

A [Go (golang)](https://golang.org) version of [paste table](https://github.com/robballou/paste-table-util).

## Example (MySQL)

Let's say you have some MySQL output like:

    +--------------------+-----+
    | information_schema | 123 |
    | mysql              | 234 |
    | performance_schema | 345 |
    | sys                | 456 |
    +--------------------+-----+

To get just the first column as a list:

    cat mysql.txt | paste-table-go

To get the second column as a list:

    cat mysql.txt | paste-table-go --column=1

To get the second column, one record per line:

    cat mysql.txt | paste-table-go --column=1 --newline

To convert to CSV:

    cat mysql.txt | paste-table-go --column=-1

## Example (CSV)

You can also parse delimited files:

    information_schema,123
    mysql,234
    performance_schema,345
    sys,456

To get just the first column as a list:

    cat mysql.txt | paste-table-go --delimiter=','

To get the second column as a list:

    cat mysql.txt | paste-table-go --delimiter=',' --column=1

To get the second column, one record per line:

    cat mysql.txt | paste-table-go --delimiter=',' --column=1 --newline

## Install (Advanced)

This project requires Go to be installed to build the binary for your system.

You can either run `bash install.sh` or run the steps manually:

    ORIGINAL_GOPATH=$GOPATH
    GOPATH=$(pwd)

    go get github.com/mkideal/cli
    go build

    GOPATH=$ORIGINAL_GOPATH

You can then add the `paste-table-go` binary to your to `$PATH` or symlink it somewhere useful.
