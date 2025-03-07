<!-- Created by mkdoc DO NOT EDIT. -->

# Examples

```sh
unitconv -from pint -to litre
```
This will show how many litres in a pint

```sh
unitconv -from chain -to mile
```
This will show how many chains in a mile

```sh
unitconv -from chain -to mile -val 80
```
This will show 80 chains in miles

```sh
unitconv -from chain -to m -val 80
```
This will show 80 chains in metres

```sh
unitconv -from chain -to m -val 80 -just-val
```
This will show 80 chains in metres\. Only the value is shown and no surrounding
explanatory text

```sh
unitconv -from chain -to m -val 80 -roughly
```
This will show 80 chains in metres\. The value is adjusted to show the nearest
multiple of 5 or 10

