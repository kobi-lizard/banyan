#!/usr/bin/env bash

SSH_KEY=$1

update(){
    SERVER_ADDR=(`cat public_ips.txt`)
    for (( j=1; j<=$1; j++))
    do
       scp -i ${SSH_KEY} config.json run.sh ips.txt $2@${SERVER_ADDR[j-1]}:~/bamboo
       ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} 'chmod 777 ~/bamboo/run.sh'
    done
}

USERNAME="ubuntu"
MAXPEERNUM=(`wc -l public_ips.txt | awk '{ print $1 }'`)

# update config.json to replicas
update $MAXPEERNUM $USERNAME
