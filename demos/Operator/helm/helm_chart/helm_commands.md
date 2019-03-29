# Helm Commands

## Useful Commands

- Create helm chart
```
helm create mychart
```

- Deploy helm chart
```
helm install ./mychart
helm install --debug --dry-run ../mychart
helm install -f myvals.yaml ./mychart
helm install --dry-run --debug --set favoriteDrink=slurm ./mychart
helm install stable/drupal --set image=my-registry/drupal:0.1.0 --set livenessProbe.exec.command=[cat,docroot/CHANGELOG.txt] --set livenessProbe.httpGet=null --dry-run
```

- Deployed helm chart
```
helm list
```

- Get all objects manifest of deployed helm chart
```
helm get manifest $NAME
```

- Delete deployed helm chart but still the helm name left
```
helm delete ${HELM_NAME}
```

- Delete deployed helm chart including helm name
```
helm del ${HELM_NAME} --purge
```