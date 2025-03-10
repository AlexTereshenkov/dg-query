#! /bin/bash -xue

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh
      
echo '--- Static checks'
echo 'Gazelle'
bazel run //:gazelle -- -mode=diff
echo 'Buildifier'
bazel run //:buildifier-check
echo 'gofmt'
output=$(bazel run @rules_go//go -- fmt .) && [[ -n "$output" ]] && echo "$output" && exit 1 || echo "No formatting required"

echo '--- Run'
bazel build //... && bazel-bin/dg-query_/dg-query dependencies --dg="examples/dg-real.json" foo.py spam.py
      
echo '--- Test unit'
bazel test --test_tag_filters=unit --test_output=all --test_summary=detailed //... --runs_per_test=20 //...

echo '--- Test integration'
bazel test --test_tag_filters=integration --test_output=all --test_summary=detailed //... 

echo '--- Test performance'
bazel test --test_tag_filters=performance --test_output=all --test_summary=detailed --runs_per_test=10 //... 

echo '--- Code coverage'
bazel coverage //... --combined_report=lcov && bazel run //:run_coverage_with_genhtml && ls coverage-html/index.html
    
echo '--- Release binary with a stamp'
bazel build //:dg-query --stamp
bazel-bin/dg-query_/dg-query | grep "Git revision"

rm .bazelrc.ci
