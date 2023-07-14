#!/usr/bin/env sh

build/couture --rate-limit=5 --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa @@fake
