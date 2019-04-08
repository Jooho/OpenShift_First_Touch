# Create Helm Chart Repo for public share

- Setup github.io for custom helm-charts
   ```
   git clone https://github.com/Jooho/helm-charts-repos.git
   
   touch index.html 
   git add index.html
   git commit -m "add index"
   git push --set-upstream origin master
   ```

- Copy helm chart package into repo folder
  ```
  cp ~/.helm/repository/local/nfs-provisioner-3.0.tgz .
  ```

- Create index file and push it
   ```
   helm repo index ../helm-charts-repos/ --url https://jooho.github.io/helm-charts-repos/

   git add .
   git commit -m "changet index.yaml"
   git push
   ```
- Add a new helm repo 
  ```
  helm repo add jooho-helm-charts https://jooho.github.io/helm-charts-repos
  ```

- Check repo and search the custom chart
  ```
  helm list
  helm repo list
  helm search nfs
  ```

- Install the custom chart  
  ```
  helm install jooho-helm-charts/nfs-provisioner --name=test
  ```
