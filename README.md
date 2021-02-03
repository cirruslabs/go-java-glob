# go-java-glob

A port of the [Java's NIO globbing functionality](https://javapapers.com/java/glob-with-java-nio/) to Golang.

## Rationale

Cirrus CI [historically used](https://github.com/cirruslabs/cirrus-ci-docs/issues/583#issuecomment-592952905) Java's NIO [PathMatcher](https://docs.oracle.com/javase/8/docs/api/java/nio/file/PathMatcher.html) for it's [`changesInclude`](https://cirrus-ci.org/guide/writing-tasks/#supported-functions) configuration function, until recently [Cirrus CLI](https://github.com/cirruslabs/cirrus-cli) was introduced, and the need for the Golang globbing implementation has arisen.

There exist a popular [doublestar](https://github.com/bmatcuk/doublestar) package for Golang, but it has subtle (and, unfortunately, deal-breaking) differences from the Java globbing. For example, in how [the `**` syntax works](https://github.com/bmatcuk/doublestar/issues/54).

Thus the original source was manually translated from the [`Globs.java`](https://github.com/openjdk/jdk/blob/3789983e89c9de252ef546a1b98a732a7d066650/src/java.base/share/classes/sun/nio/fs/Globs.java) file publicly available as a part of the [OpenJDK](https://github.com/openjdk/jdk) release.

The result can be found in the [`glob.go`](glob.go), which maps in a one-to-one fashion to the original (except for the change that works around missing `&&` (AND) operator in Golang's regexp package) and can be compared for differences in a matter of a few minutes.

## Usage

The package can be imported as follows:

```
github.com/cirruslabs/go-java-glob
```

The simplest invocation looks like this:

```go
re, err := glob.ToRegexPattern("**.go")
if err != nil {
	return err
}

if re.MatchString(path) {
	// the path has matched
}
```
