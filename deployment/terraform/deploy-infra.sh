#!/bin/bash
# #run tofu init
tofu init

#Get cluster name:
clusterName=$(echo local.name | tofu console | tr -d '"' )

#Get aws region:
awsRegion=$(echo var.aws_region | tofu console | tr -d '"' )

#export AWS_PROFILE=default
export AWS_PROFILE=cb
export AWS_REGION=${awsRegion}

##just base infra to be able to run the rest of the modules
tofu plan -out=plan.tfplan --target=module.eks --target=module.karpenter --target=helm_release.karpenter --target=helm_release.aws-load-balancer-controller --target=aws_eks_addon.aws_ebs_csi_driver
#run tofu plan
tofu apply plan.tfplan

#run tofu plan saving the plan
tofu plan -out=plan.tfplan

# run tofu apply using the plan
tofu apply plan.tfplan

#update the kubectl file
aws eks update-kubeconfig --name ${clusterName} --alias ${clusterName}

# some prints
grafana_pass=$(kubectl --namespace prometheus get secrets prometheus-grafana -o jsonpath="{.data.admin-password}" --context ${clusterName} | base64 -d)
echo "grafana password: ${grafana_pass}"


