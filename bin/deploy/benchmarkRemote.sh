#!/bin/bash

# Example usage:
# bash benchmarkRemote.sh banyan 0 1 1
# bash benchmarkRemote.sh icc 0 1 1
# bash benchmarkRemote.sh hotstuff 0 1 1
# bash benchmarkRemote.sh streamlet 0 1 1

#---------------------------------------
LOG_DIR="logs"
DURATION=120

#---------------------------------------
# SSH key with which the aws servers are configured
# TODO adapt the aws/settings.json file: set the same value as here.
SSH_KEY="/local/home/yvonlanthen/.ssh/aws_global"

#---------------------------------------
ALGORITHM=$1                # Algorithm run
N_CLIENT_THREADS=$2         # (deprecated) Number of txs sent to replicas (we now steer this through block size in config.json)
NEW_AWS=$3                  # Set to 1 if its the first time you run this aws config
BUILD=$4                    # Set to 1 if binaries should be built
TIMESTAMP=$(date +%H:%M:%S)
LOG_FILE="${1}_${2}_${TIMESTAMP}.txt"
#---------------------------------------

echo  "--- Start"
mkdir $LOG_DIR

if [ $NEW_AWS -eq 1 ] 
then
    echo "Get IPs"
    cd ../../aws

    fab ipservers > ../bin/deploy/public_ips.txt
    fab ipservers > ../bin/deploy/ips.txt
    echo "servers:"
    cat '../bin/deploy/public_ips.txt'
    
    cd ../bin/deploy

fi

if [ $BUILD -eq 1 ] 
then
    echo "--- Build"
    chmod +777 build.sh
    ./build.sh

    echo "--- Setup Servers"
    chmod +777 deploy.sh
    ./deploy.sh $SSH_KEY
fi


echo "--- Update Configuration"
chmod +777 update_conf.sh
./update_conf.sh $SSH_KEY

echo "--- Start Servers"
chmod +777 start.sh
./start.sh $ALGORITHM $SSH_KEY

echo "--- wait "
sleep 5


probe(){
    SERVER_ADDR=(`cat public_ips.txt`)
    for (( j=1; j<=$1; j++))
    do 
        HTTP_PORT=$(( 8069 + ${j-1} ))
        echo $HTTP_PORT
        echo '' >> $LOG_DIR/$LOG_FILE
        curl -i -H "Accept: text/plain" -H "Content-Type: text/plain" -X GET http://${SERVER_ADDR[j-1]}:${HTTP_PORT}/query >> $LOG_DIR/$LOG_FILE
        echo "query"
    done
}

echo "--- Kickoff Servers"

SERVER_ADDR=(`cat public_ips.txt`)
curl -i -H "Accept: text/plain" -H "Content-Type: text/plain" -X GET http://${SERVER_ADDR[0]}:8070/query
# MAXPEERNUM=(`wc -l public_ips.txt | awk '{ print $1 }'`)
# probe $MAXPEERNUM


sleep $DURATION

echo "--- Probe Servers"

MAXPEERNUM=(`wc -l public_ips.txt | awk '{ print $1 }'`)
probe $MAXPEERNUM

echo "--- Log output for ${ALGORITHM} with ${N_CLIENTS} clients:"
cat $LOG_DIR/$LOG_FILE

echo "--- Shutdown"

./pkill.sh $SSH_KEY

