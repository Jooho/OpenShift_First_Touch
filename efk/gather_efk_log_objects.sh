
# Temp directory
export temp_dir=log_temp_dir
mkdir ${temp_dir}
cd ./${temp_dir}

# Save all objects
for i in $(oc get all,secret,sa,oauthclient -o name)
do
 echo "Save ${i} => $(echo ${i}|sed 's/\//_/g' ).yaml"

 oc get $i -o yaml > $(echo ${i}|sed 's/\//_/g').yaml 
done

# Kibana pod name
kibana_pod=$(oc get pod -n logging|grep kibana|awk '{print $1}')

# Kibana Log
echo "Save Kibana log to ${kibana_pod}_kibana.log"
oc logs ${kibana_pod} -c kibana > ${kibana_pod}_kibana.log

# Kibana-proxy Log
echo "Save Kibana Proxy log to ${kibana_pod}_kibana_proxy.log"
oc logs ${kibana_pod} -c kibana-proxy > ${kibana_pod}_kibana-proxy.log

# ES/FluentD Log
for i in $(oc get pod -n logging|egrep -v "kibana|NAME" |awk '{print $1}');
do
 echo "Save ${i} to ${i}.log"
 oc logs ${i} > ${i}.log
done

# ES health/indices
for i in $(oc get pod -n logging|grep es|awk '{print $1}');
do
 touch ./${i}_es.output
 echo "oc exec ${i} -- curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key https://localhost:9200/_cat/health?v" >> ${i}_es.output
 oc exec ${i} -- curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key https://localhost:9200/_cat/health?v >> ${i}_es.output
 echo "" >>  ${i}_es.output


 echo "oc exec ${i} -- curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key https://logging-es:9200/_cat/indices?v" >> ${i}_es.output
 oc exec ${i} -- curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key https://logging-es:9200/_cat/indices?v >> ${i}_es.output
 echo "" >>  ${i}_es.output

 echo "oc exec ${i} -- du -sch /elasticsearch/persistent" >> ${i}_es.output
 oc exec ${i} -- du -sch /elasticsearch/persistent >> ${i}_es.output
 echo "" >>  ${i}_es.output

 echo "oc exec ${i} -- df -h" >> ${i}_es.output
 oc exec ${i} -- df -h >> ${i}_es.output
 echo "" >>  ${i}_es.output

done

tar cvf ../log.tar ./*
cd ..

echo "Please delete ${temp_dir} manaually"
echo "Please upload log.tar file"

