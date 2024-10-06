#!/usr/bin/env bash

SSH_KEY=$1

# add ssh-key
# add_ssh_key(){
# 	SERVER_ADDR=(`cat public_ips.txt`)
#     echo "Add your local ssh public key into all nodes"
#     for (( j=1; j<=$1; j++ ))
#     do
#             addkey ${SERVER_ADDR[j-1]} $2 $3
# 	    wait
#     done
# }

distribute(){
    SERVER_ADDR=(`cat public_ips.txt`)
    for (( j=1; j<=$1; j++))
    do 
       ssh-keyscan ${SERVER_ADDR[j-1]} >> ~/.ssh/known_hosts
       ssh -i ${SSH_KEY} -t $2@${SERVER_ADDR[j-1]} mkdir bamboo
       echo -e "---- upload replica ${j}: $2@${SERVER_ADDR[j-1]} \n ----"
       scp -i ${SSH_KEY} server ips.txt $2@${SERVER_ADDR[j-1]}:~/bamboo
    done
}

USERNAME='ubuntu'
MAXPEERNUM=(`wc -l public_ips.txt | awk '{ print $1 }'`)

# distribute files
distribute $MAXPEERNUM $USERNAME
