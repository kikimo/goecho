#!/bin/bash

for i in $(seq 10 40 500); do
	./client $i
done
