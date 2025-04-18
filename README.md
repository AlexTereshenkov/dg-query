# Documentation

`dg-query` is a command-line utility program to query dependency graph of a codebase. It operates on the [adjacency lists](https://en.wikipedia.org/wiki/Adjacency_list) data stored as a JSON file. Note that this is not a graph library; for advanced graph manipulation and analysis, [`networkx`](https://networkx.org/) would be a great choice. 

The tool is written in Go, built with Bazel, and runs its CI in Buildkite.

## Dependency graph

The `dg-query` is intended to be used to work primarily with the snapshot of the dependency graph of a codebase, with every build target mapped to an array of its direct dependencies e.g.

```json
{
    "src/app/main.py": [
        "src/shared/utils.py", 
        "src/shared/process.py"
    ],
    "src/lib/project.py": [
        "src/lib/internal/billing/invoice.py", 
        "src/lib/internal/rpc/process.py"
    ]
}
```

If you don't have a sophisticated build system, this kind of data should still be possible to obtain by hacking together some scripts to statically analyze your import statements iteratively building the adjacency lists. 

If you use an artifact based build system such as [Pants](https://www.pantsbuild.org/) or [Bazel](https://bazel.build/), exporting the dependency graph is trivial, but may still require some post-processing such as omitting irrelevant subsets of the graph or renaming some build targets for clarity.

Pants (see [command-line reference](https://www.pantsbuild.org/stable/reference/goals/dependencies)):

```shell
$ pants dependencies --format=json :: > dg.json
```

Bazel (see [Querying dependency graph of a Bazel project](https://alextereshenkov.github.io/querying-dependency-graph-bazel.html)):

```shell
$ bazel query 'deps(//...)' --output=graph --noimplicit_deps > graph.dot
# convert dot file to JSON
```

Build systems allow exporting data about the reverse dependencies (aka dependents), but this is not required for the `dg-query` as it operates solely on the dependencies lists.

## Features

The tool serves as a faster way to query the dependency graph of a project. This is because the targets and their relationships are materialized so there's no need for any kind of runtime evaluation. The downside is that you need to re-export your dependency graph data every time a change that leads to changes in the dependency graph is made.

`dg-query` has a ton of functionality grouped under individual commands:

### `dependencies` (`deps`)
Identify dependencies of given node(s), optionally transitively (`--transitive`).

Options:
* specify depth when searching for dependencies transitively (`--depth`)
* include the build target itself in the output (`--reflexive`)

### `dependents` (`rdeps`)
Identify dependents of given node(s), optionally transitively (`--transitive`).

Options:
* specify depth when searching for dependents transitively (`--depth`)
* include the build target itself in the output (`--reflexive`)

### `roots`
Get nodes that no other node depends on. The roots are also known as sources. 

### `leaves`
Get nodes that have no dependencies. The leaves are also known as sinks. 

### `paths`
Get paths between individual targets.

Options:
* `--from` and `--to` targets to find paths between
* `--n` limits the number of paths returned (helpful with a large graph)

### `cycles`
Find [cycles](https://en.wikipedia.org/wiki/Cycle_(graph_theory)) in the dependency graph.
This is useful when you use a build system that doesn't tolerate cycles and you want to
get a list of all of them at once.

### `components`
Find [components](https://en.wikipedia.org/wiki/Component_(graph_theory)) in the dependency graph.
This is useful when you want to find out how well your repository is separated in terms of independent
modules or projects.

### `subgraph`
Extract [subgraph](https://en.wikipedia.org/wiki/Glossary_of_graph_theory#subgraph) out of the dependency graph.
This is useful when you want to visualize a subset of the dependency graph or study it closer.

### `simplify`
Simplify the dependency graph applying certain techniques. 
Currently support: [Transitive reduction](https://en.wikipedia.org/wiki/Transitive_reduction). 
This is useful when you want to make graph visualization less cluttered or to compact a very large graph.

### `metrics`
Get dependency graph related metrics. A dependency graph (`--dg`) may be used,
or a reverse dependency graph (`--rdg`) may be used, if you have one.

Metrics:
* dependency count (optionally, transitively)
* dependent count (optionally, transitively)
* [connected components](https://en.wikipedia.org/wiki/Component_(graph_theory)) count (few components suggests a very tight graph)
