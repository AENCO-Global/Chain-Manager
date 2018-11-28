#!/usr/bin/env bash

#!/bin/bash
# ---------------------------------------------------
echo "--=== Incoming Paramters (This script should be reusable) ===--"
echo "[P1] Version Number is :$1 "
echo "[P2] Target Server is :$2 "
echo "[P3] Target Folder is :$3 "
echo "[P4] Target Folder is :$4 "
echo "[P5] Target Folder is :$5 "
echo "---------------------------------------"

echo "--=== Modify Version Information ===--"
echo "Version $1" > ./version.info

ls  -l
echo "--------------------------------------"
./setupBuild.sh
echo "--======== Build the Application =======--"
./build.sh

echo "--=== Transfer files to remote Server ===--"
echo "rsync -avzhe ssh  --rsync-path="""rsync""" ./ jenkins@$2:$3"""
ssh -o "StrictHostKeyChecking=no" -p 22 $2 "pwd"
if [ $2 == 'build-server.aencoin.com' ] ; then
    rsync -avzhe ssh  --rsync-path="rsync" ./aen.agent jenkins@$2:$3/_build/bin
    rsync -avzhe ssh  --rsync-path="rsync" ./resources/config-agent.ini jenkins@$2:$3/resources
    ssh -p 22 $2 "sed -i \"s/{id}/$4/g\" $3/resources/config-agent.ini"
    ssh -p 22 $2 "sed -i \"s/{version}/$1/g\" $3/resources/config-agent.ini"
else
    ssh -p 22 $2 "mkdir -p $3/services"
    rsync -avzhe ssh  --rsync-path="rsync" ./aen.agent jenkins@$2:$3/chain
    rsync -avzhe ssh  --rsync-path="rsync" ./resources/config-agent.ini jenkins@$2:$3/resources
    rsync -avzhe ssh  --rsync-path="rsync" ./services/* jenkins@$2:$3/services
    rsync -avzhe ssh  --rsync-path="rsync" ./devops.sh jenkins@$2:$3/

    echo "Updating the configuration file"
    ssh -p 22 $2 "sed -i \"s/{id}/$4/g\" $3/resources/config-agent.ini"
    ssh -p 22 $2 "sed -i \"s/{version}/$1/g\" $3/resources/config-agent.ini"
    echo "ReStarting the Service"
    ssh -p 22 $2 "sudo systemctl daemon-reload &"
    ssh -p 22 $2 "sudo systemctl restart aen-agent &"
fi

echo "----====== Verify Deployments-List from Remote ======----"
ssh -p 22 $2 "cat $3/resources/config-agent.ini"
ssh -p 22 $2 "systemctl status | grep aen.agent | grep -v grep "
echo "---------------------------------------------------------"
ssh -p 22 $2 "ls -al $3"
echo "---------------------------------------------------------"

echo "--=== Version Deployed is [$1] The following output from version.info ===--"
echo "------------The-End-------------------------------------------------------"
