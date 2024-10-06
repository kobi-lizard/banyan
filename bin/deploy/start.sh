#!/usr/bin/env bash
ALGORITHM=$1
SSH_KEY=$2

./pkill.sh $SSH_KEY

start(){
    SERVER_ADDR=(`cat public_ips.txt`)
    for (( j=1; j<=$1; j++))
    do
      ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} "cd ~/bamboo ; nohup ./run.sh ${j} ${ALGORITHM}"
      sleep 0.1
      echo replica ${j} is launched!
    done
}

USERNAME="ubuntu"
MAXPEERNUM=(`wc -l public_ips.txt | awk '{ print $1 }'`)

# update config.json to replicas
start $MAXPEERNUM $USERNAME
