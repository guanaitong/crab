kubectl apply -f crab-rbac.yaml

kubectl \
-n kube-system \
describe secret \
$(kubectl -n kube-system get secret | grep crab | awk '{print $1}')