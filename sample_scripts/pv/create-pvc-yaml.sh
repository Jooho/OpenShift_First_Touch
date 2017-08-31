




VOL_NAME=ose09
PVC_ACCESS_MODE=ReadWriteMany
LVM_VOL_SIZE=500


cat << EOF > ./${VOL_NAME}.yaml
apiVersion: "v1"
kind: "PersistentVolumeClaim"
metadata:
  name: ${VOL_NAME}
spec:
  accessModes:
    - ${PVC_ACCESS_MODE}
  resources:
    requests:
      storage: ${LVM_VOL_SIZE}Gi
EOF
  echo "Created def file for ${VOL_NAME}"


