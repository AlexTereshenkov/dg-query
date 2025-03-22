#! /bin/bash 
set -xueo pipefail

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh
      
echo '--- Static checks'
echo 'Gazelle'
bazel run //:gazelle -- -mode=diff
echo 'Buildifier'
bazel run //:buildifier-check
bazel run //:buildifier
echo 'gofmt'
output=$(bazel run @rules_go//go -- fmt .) && [[ -n "$output" ]] && echo "$output" && exit 1 || echo "No formatting required"

echo '--- Run'
bazel build //... && bazel-bin/dg-query_/dg-query dependencies --dg="examples/dg-real.json" foo.py spam.py
      
echo '--- Test unit'
bazel test --test_tag_filters=unit --test_output=all --test_summary=detailed //... --runs_per_test=20 //... | sed 's/^=== //' | sed 's/^--- //'

echo '--- Test integration'
bazel test --test_tag_filters=integration --test_output=all --test_summary=detailed //... | sed 's/^=== //' | sed 's/^--- //'

echo '--- Test performance'
bazel test --test_tag_filters=performance --test_output=all --test_summary=detailed --runs_per_test=10 //... | sed 's/^=== //' | sed 's/^--- //'

echo '--- Code coverage'
bazel coverage //... --combined_report=lcov | sed 's/^=== //' | sed 's/^--- //' && bazel run //:run_coverage_with_genhtml && ls coverage-html/index.html
    
echo '--- Release binary with a stamp'
bazel build //:dg-query --stamp
bazel-bin/dg-query_/dg-query | grep "Git revision"

echo '--- Running a compiled CGo binary'
bazel run //platforms/compiled:cgo-compiled

rm .bazelrc.ci

echo '--- Finish'
