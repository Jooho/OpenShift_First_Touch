Debugging - Gathering EFK Object
-------------------------------

In order to debug EFK stack, gathering information of EFK Objects is necessary.

[This script](./gather_efk_log_objects.sh) help to get following objects or logs:
- Objects: all,secret,sa,oauthclient

- Kibana/Kibana-Proxy Log

- ES/FluentD Log

- ES health/indices

- ES storage size


Example)
```
$ ./gather_efk_log_objects.sh
..
..
Please delete ${temp_dir} manaually
Please upload log.tar file
```


