#!/usr/bin/env bash

set -eox pipefail

function test-local() {
  set +x

  docker-compose -f make/docker-compose-dev.yml down --remove-orphans && \
  docker-compose -f make/docker-compose-dev.yml up --build -d

  export CONFIG_FILE=$(pwd)/config/config_dev.toml

  test
}

function test() {
  set +x
  set -e # exit if a command fails

  PKG_LIST=$(go list ./... | grep -v /vendor/)
  COVERAGE_DIR="${COVERAGE_DIR:-.coverage}"

  echo "tyding mods"
  go mod tidy

  echo "Formating and linting ..."
  go fmt $PKG_LIST

  # for CI
  GOPATH=$1
  if [ -z "$GOPATH" ]; then
      golint -set_exit_status $PKG_LIST
  else
      $GOPATH/bin/golint -set_exit_status $PKG_LIST
  fi

  echo "Running tests ..."

  # stop tests at first test fail
  TFAILMARKER="FAIL:"
  go test $PKG_LIST -v -count=1 -p=1 | { IFS=''; while read line; do
      echo "$line"
      if [ -z "$line" ]; then
          continue
      fi

      if [ -z "${line##*$TFAILMARKER*}" ] ; then
          echo "🚨 Test FAIL match, exit"
          exit 1
      fi
  done }

  echo "Running code coverage ..."

  # Create the coverage files directory
  mkdir -p "$COVERAGE_DIR";

  # Create a coverage file for each package
  # test minim coverage
  MINCOVERAGE=75

  for package in $PKG_LIST; do
      pkgcov=$(go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package")
      # echo $pkgcov

      case $pkgcov in
          *coverage:*)
              pcoverage=$(echo $pkgcov| grep "coverage" | sed -E "s/.*coverage: ([0-9]*\.[0-9]+)\% of statements/\1/g")
              # echo "coverage: $pcoverage% of $package"

              if [ $(echo ${pcoverage%%.*}) -lt $MINCOVERAGE ] ; then
                  echo "🚨 Test coverage of $package is $pcoverage%"
                  echo "FAIL"
                  exit 1
              else
                  echo "🟢 Test coverage of $package is $pcoverage%"
              fi
              ;;
          *)
              echo "➖ No tests for $package"
              ;;
      esac
  done

  # Merge the coverage profile files
  echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov
  for fcov in "${COVERAGE_DIR}"/*.cov
  do
      if [ $fcov != "${COVERAGE_DIR}/coverage.cov" ]; then
          tail -q -n +2 $fcov >> "${COVERAGE_DIR}"/coverage.cov
      fi
  done


  # global code coverage
  pcoverage=$(go tool cover -func="${COVERAGE_DIR}"/coverage.cov | grep 'total:' | sed -E "s/^total:.*\(statements\)[[:space:]]*([0-9]*\.[0-9]+)\%.*/\1/g")
  echo "coverage: $pcoverage% of project"

  if [ $(echo ${pcoverage%%.*}) -lt $MINCOVERAGE ] ; then
      echo "🚨 Test coverage of project is $pcoverage%"
      echo "FAIL"
      exit 1
  else
      echo "🟢 Test coverage of project is $pcoverage%"
  fi
}