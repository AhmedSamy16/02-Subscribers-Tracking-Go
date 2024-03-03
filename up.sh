#!/bin/bash

cd migrations

goose postgres postgres://postgres:@localhost:5432/subscribers up