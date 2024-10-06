#!/usr/bin/env bash

N_CLIENT_THREADS=$1
SSH_KEY=$2

start(){
    CLIENT_ADDR=(`cat clients.txt`)
    for (( j=1; j<=$1; j++))
    do
      ssh -i ${SSH_KEY} -t $2@${CLIENT_ADDR[j-1]} "cd ~/bamboo ; nohup ./runClient.sh ${N_CLIENT_THREADS}"
      sleep 0.1
      echo client ${j} is launched!
    done
}

USERNAME="ubuntu"
MAXPEERNUM=(`wc -l clients.txt | awk '{ print $1 }'`)

# update config.json to replicas
start $MAXPEERNUM $USERNAME