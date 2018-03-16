#!/bin/bash -e
cd $(dirname $0)

function announce
{
  echo
  echo $@
}

PATH=$HOME/gopath/bin:$GOPATH/bin:$PATH

. go-get.sh

# delete artefacts from previous build (if any)
mkdir -p reports
rm -f reports/*.out reports/*.html */*.txt demo/*_sql.go

### Collection Types ###
# these generated files hardly ever need to change (see github.com/rickb777/runtemplate to do so)
[ -f schema/type_set.go ]   || runtemplate -tpl simple/set.tpl  -output schema/type_set.go  Type=Type   Comparable:true Ordered:false Numeric:false
[ -f util/int64_set.go ]    || runtemplate -tpl simple/set.tpl  -output util/int64_set.go   Type=int64  Comparable:true Ordered:true  Numeric:true
[ -f util/string_list.go ]  || runtemplate -tpl simple/list.tpl -output util/string_list.go Type=string Comparable:true Ordered:true  Numeric:false
[ -f util/string_set.go ]   || runtemplate -tpl simple/set.tpl  -output util/string_set.go  Type=string Comparable:true Ordered:true  Numeric:false

### Build Phase 1 ###

cd sqlgen
go install .

for d in code output parse; do
  announce sqlgen/$d
  go test ./$1 -covermode=count -coverprofile=../reports/sqlgen-$d.out ./$d
  go tool cover -html=../reports/sqlgen-$d.out -o ../reports/sqlgen-$d.html
  [ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=../reports/sqlgen-$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
done

cd ..

### Build Phase 2 ###

go install ./...

for d in constraint require schema sqlgen where; do
  announce ./$d
  go test $1 -covermode=count -coverprofile=reports/$d.out ./$d
  go tool cover -html=reports/$d.out -o reports/$d.html
  [ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
done

#echo .
#go test . -covermode=count -coverprofile=dot.out .
#go tool cover -func=dot.out
#[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN

### Demo ###

announce demo
cd demo
./build.sh sqlite mysql postgres

cd ..
echo
echo go vet -shadow=true -composites=false ./...
go vet -shadow=true -composites=false ./...
