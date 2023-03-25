#!/bin/bash

SLAVE1_NODE_NAME=$(kubectl get nodes -o wide | awk '$6 == "192.168.50.11" {print $1}')
kubectl label node $SLAVE1_NODE_NAME node-role.kubernetes.io/worker=""
