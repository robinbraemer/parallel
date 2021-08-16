# parallel
`parallel` is a tiny utility package for running functions in parallel with error handling.
You can very easily orchestrate highly complex function executions that run in parallel or ordered fashion,
combine both and get a single standard Go error back that supports functions from Go's errors package like
`errors.Is` and `errors.As`.

## Install
```bash
go get -u github.com/robinbraemer/parallel
```

## Example

Imagine you need to run multiple functions in parallel
but some of them in order. It would be quiet code heavy
to do so in a clean manner.

Luckily you can use the `parallel` package to de-complicate this:

```go
var fns [][]Fn
// ...
err := Ordered(
    Parallel( // A
        Parallel(fns[0]...), // 1
        Parallel(fns[1]...), // 2
        Parallel(fns[2]...), // 3
    ), 
    Parallel(fns[3]...),  // B
    Ordered(
        Parallel(fns[5]...), // 4
        Parallel(fns[6]...), // 5
    ),  // C
).Do()
```

This is what happens:
- Ordered runs A, B, C in unparalleled order and returns the first error encountered
- `A` runs 3 slices of functions in parallel and block until all are finished
- after `A` is done `B` runs functions in parallel
- after `B` is done `C` runs functions in parallel
- Ordered returns

That means `1-3` must complete before `A` returns and `B` is run and so forth.