#!/usr/bin/env bash

SSH_KEY=$1

distribute(){
    SERVER_ADDR=(`cat clients.txt`)
    for (( j=1; j<=$1; j++))
    do
       ssh-keyscan ${SERVER_ADDR[j-1]} >> ~/.ssh/known_hosts
       ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} mkdir bamboo
       echo -e "---- upload client ${j}: $2@${SERVER_ADDR[j-1]} \n ----"
       scp -i ${SSH_KEY} client ips.txt config.json runClient.sh closeClient.sh $2@${SERVER_ADDR[j-1]}:~/bamboo
       ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} "chmod 777 ~/bamboo/runClient.sh"
       ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} "chmod 777 ~/bamboo/closeClient.sh"
       wait
    done
}

USERNAME='ubuntu'
MAXPEERNUM=(`wc -l clients.txt | awk '{ print $1 }'`)

# distribute files
distribute $MAXPEERNUM $USERNAME
