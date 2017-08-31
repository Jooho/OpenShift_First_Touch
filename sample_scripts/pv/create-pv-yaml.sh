
VOL_NAME=pv-infra-ose01g0001
LVM_VOL_SIZE=1
PV_ACCESS_MODE=ReadWriteMany
NFS_MOUNT_POINT=/exports
LVM_VOL_NAME=lv_001
NFS_SERVER=infra.example.com
PV_RECLAIM_POLICY=Retain

cat << EOF > ${VOL_NAME}.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: "${VOL_NAME}"
spec:
  accessModes:
  - ${PV_ACCESS_MODE}
  capacity:
    storage: ${LVM_VOL_SIZE}Gi
  nfs:
    path: ${NFS_MOUNT_POINT}/${LVM_VOL_NAME}
    server: ${NFS_SERVER}
  persistentVolumeReclaimPolicy: ${PV_RECLAIM_POLICY}
EOF
  echo "Created def file for ${VOL_NAME}"


