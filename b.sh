#!/bin/bash
docker start -a build-gocqplg
./cqcfg . -c
rm -rf ~/coolq/dev/me.cqp.molin.secretmaster
mkdir -p ~/coolq/dev/me.cqp.molin.secretmaster
cp app.dll ~/coolq/dev/me.cqp.molin.secretmaster/
cp app.json ~/coolq/dev/me.cqp.molin.secretmaster/

