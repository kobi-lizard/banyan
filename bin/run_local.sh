#!/usr/bin/env bash

go build ../server/
./server -sim=true -log_level=debug -algorithm=banyan