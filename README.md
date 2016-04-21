## tree - scapegoat tree implementation
Scapegoat tree is height balanced binary search tree

It could be parametrized by single `alfa` number in range `[0.5-1)`.

Assuming that `alfa` is constant, both insert&delete and search operations have `O(lg n)` complexity.

For _bigger_ `alfa` values, insertion will be less likely to trigger rebalance of (part of) the tree.
For _smaller_ values, searches will be faster, but at insertion cost.

# usage
Copy code and change key type (`byte` in example) to desired type.
You need also look into `cmp()` function.

If you want to use additional value, just use struct with more fields and skip those fields during comparison.

And, please, don't comply about generics :)


# development stage
For now, structure is as usable as C++' STL's set/map - ordered container with logarithm time access.

I plan to add support for _augmented_ operations, like access to _k_-th element, sum in range, etc.

# TODO
- [ ] document all Exported functions;
- [ ] provide better rule of thumb for `alfa` parameter, for now please use 0.65 or google;
