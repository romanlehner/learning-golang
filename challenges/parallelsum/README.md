# Parallel Sum

## Problem statement

Create a function Sum() that takes a slice of integers and a go channel and sends the sum of integers into the channel. Create a test that runs the Sum() function as go routine and collect all the sums through the go channel.

Run tests with:

```bash
cd challenges
go test ../challenges/parallelsum/
```