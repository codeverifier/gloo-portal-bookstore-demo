#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

Provision() {
  echo "------------------------------------------------------------"
  echo "Deploying application"
  echo "------------------------------------------------------------"
  echo ""
  kubectl apply -f $DIR/deployment/mongo.yaml -n apps
  kubectl apply -f $DIR/deployment/bookstore.yaml -n apps
  kubectl apply -f $DIR/deployment/bookstore-upstream.yaml -n apps-configuration

  echo ""
  echo "------------------------------------------------------------"
  echo "Injecting assets for publishing APIs and portal"
  echo "------------------------------------------------------------"
  echo ""
  kubectl apply -f $DIR/api-publish/bookstore-route.yaml
  kubectl apply -f $DIR/api-publish/bookstore-schema.yaml
  kubectl apply -f $DIR/api-publish/bookstore-product.yaml
  envsubst < <(cat $DIR/api-publish/bookstore-env.yaml) | kubectl apply -f -
  envsubst < <(cat $DIR/api-publish/bookstore-portal.yaml) | kubectl apply -f -

  echo ""
  echo "---- Generating a user/group for testing ... ----"
  pass=$(htpasswd -bnBC 10 "" Passwd00 | tr -d ':\n')
  # Store the hash as a Kubernetes Secret
  kubectl create secret generic dev1-password -n apps-configuration \
    --type=opaque \
    --from-literal=password=$pass 
  kubectl apply -f $DIR/api-publish/acl/dev-group.yaml
  kubectl apply -f $DIR/api-publish/acl/dev-user.yaml 
}

Delete() {
  echo "Cleaning up ..."

  kubectl delete -f $DIR/api-publish/bookstore-route.yaml
  kubectl delete -f $DIR/api-publish/bookstore-schema.yaml
  kubectl delete -f $DIR/api-publish/bookstore-product.yaml
  envsubst < <(cat $DIR/api-publish/bookstore-env.yaml) | kubectl delete -f -
  envsubst < <(cat $DIR/api-publish/bookstore-portal.yaml) | kubectl delete -f -

  kubectl delete secret dev1-password -n gloo-portal
  kubectl delete -f $DIR/api-publish/acl/dev-group.yaml
  kubectl delete -f $DIR/api-publish/acl/dev-user.yaml

  kubectl delete -f $DIR/deployment/mongo.yaml -n apps
  kubectl delete -f $DIR/deployment/bookstore.yaml -n apps
  kubectl delete -f $DIR/deployment/bookstore-upstream.yaml -n apps-configuration
}

shift $((OPTIND-1))
subcommand=$1; shift
case "$subcommand" in
    prov )
        Provision
    ;;
    del )
        Delete
    ;;
    * ) # Invalid subcommand
        if [ ! -z $subcommand ]; then
            echo "Invalid subcommand: $subcommand"
        fi
        exit 1
    ;;
esac