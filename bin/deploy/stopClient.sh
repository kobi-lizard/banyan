#!/usr/bin/env bash

SSH_KEY=$1

stop(){
    SERVER_ADDR=(`cat clients.txt`)
    for (( j=1; j<=$1; j++))
    do
      ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} "cd ~/bamboo ; nohup ./closeClient.sh"
      sleep 0.1
      echo client ${j} is closed!
    done
}

USERNAME="ubuntu"
MAXPEERNUM=(`wc -l clients.txt | awk '{ print $1 }'`)

# update config.json to replicas
stop $MAXPEERNUM $USERNAME