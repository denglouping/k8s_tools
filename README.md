**k8s command line tools**
-
**makeConfig**

introduce:build a kubeconfig file with specified parameters at home dir

parameters:

1.the path of kubeconfig file to use
     
2.the permission for new kubeconfig file.

3.the namespace for new kubeconfig file.

4.the serviceaccount for binding permisson

example:

    [root@test ~]# ./makeConfig -k=.kube/config -s world -p admin -n hello
    
    [root@test ~]# kubectl get sa  -n hello world -o yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      creationTimestamp: "2020-04-19T14:51:28Z"
      name: world
      namespace: hello
      resourceVersion: "1344403"
      selfLink: /api/v1/namespaces/hello/serviceaccounts/world
      uid: 2c068021-a5ef-4f20-93e2-78ceddb0d0f3
    secrets:
    - name: world-token-ht5cl
    
    [root@test ~]# kubectl get rolebinding -n hello world-hello-admin -o yaml
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      creationTimestamp: "2020-04-19T14:51:28Z"
      name: world-hello-admin
      namespace: hello
      resourceVersion: "1344404"
      selfLink: /apis/rbac.authorization.k8s.io/v1/namespaces/hello/rolebindings/world-hello-admin
      uid: 7f1033df-554b-4511-967c-d7351ae4a5fe
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: admin
    subjects:
    - kind: ServiceAccount
      name: world
      namespace: hello
    You have new mail in /var/spool/mail/root
    
    [root@test ~]# cat world-hello-admin.config 
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUR2akNDQXFhZ0F3SUJBZ0lVUzVlaEZNK3Z4bm5VK1lJMk1PMzNHV05BUzZJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pURUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFVcHBibWN4RURBT0JnTlZCQWNUQjBKbAphVXBwYm1jeEREQUtCZ05WQkFvVEEyczRjekVQTUEwR0ExVUVDeE1HVTNsemRHVnRNUk13RVFZRFZRUURFd3ByCmRXSmxjbTVsZEdWek1CNFhEVEU1TURrek1EQXpNVE13TUZvWERUSTBNRGt5T0RBek1UTXdNRm93WlRFTE1Ba0cKQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFVcHBibWN4RURBT0JnTlZCQWNUQjBKbGFVcHBibWN4RERBSwpCZ05WQkFvVEEyczRjekVQTUEwR0ExVUVDeE1HVTNsemRHVnRNUk13RVFZRFZRUURFd3ByZFdKbGNtNWxkR1Z6Ck1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBeGw3UGRtRDJ5RlkvVEtSY09jbUYKVy9WcEs3RmhVM2YwYmprNURSRUNBNlhxRDZMQ0V1VVhGV20yS1YwcWEzc0Jxc2JvZUkrazcyWFc5R2t5NkNYUAoramtIQkxuUE1wcmdLRmtvdXY4MUQwQlVwZ2h1WlV2d3VLNGEweEY4dkpxN3BXY3l1OEVuTUMwNm1WbGU4QytECk43Qm5SeFNVVW8vVVowY1BHbTBGRE15dEZ3RENIdFlMUWN0ME5BUzV1OVVTY2p6MFdLSm9XbSs0RHl6TVhrM08KeUhJNGV0Nzlrb2xTQmJyb25sNk14SC9TRXdFTGM1aTcxaGNHWHlDWFZxZjFDTUdyODlBalNjaVRpc1d1Rk5SWgpYRkkxei9xeUkzZVRNZi9Fc0ZqQnVYNUJ3eXVQSkxuY05YS3FpVW8xdEQvZlA5bDUvZ00rQlJEOExhYm1tMVg5CkZRSURBUUFCbzJZd1pEQU9CZ05WSFE4QkFmOEVCQU1DQVFZd0VnWURWUjBUQVFIL0JBZ3dCZ0VCL3dJQkFqQWQKQmdOVkhRNEVGZ1FVL3BkcnI4RnhBZHdSRHFGZXdNdUVjaStkUVRJd0h3WURWUjBqQkJnd0ZvQVUvcGRycjhGeApBZHdSRHFGZXdNdUVjaStkUVRJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFEMWRPclVrZ3JUWDUyNklvZDFRCkN5RUNSQklHVW1YMUkxcS9JTm1oK3NNQlRlSUtyUEFMV0NiRTR5Wmd4b24xQXVXVDZ3ZCtUTUE2a1puSFM3cC8KUERzSVlETC9ULzJkeUR6SmNDK2dPZG1QbWpOMHI2V1dZNTlDUU9OUVZzZzMxbFB0Tks2YkhDYlR2K2ovU09yMQpjanJJY0ZnSWRLQld6QjFsbFU5dkhKZnFBdFp5U29xOTh4UG5BYjJaTExnKzdSb0hRdzVnRkFSSWZrNnJqN3FvCnZYUEp0T1VHSGZoMVV3T1dCaVpuWWl1Rk44cDJwcjlCYkVaQ1QxTVpJSzhaMFNUN3VpR09MZjhjZ0IwaXZGTGoKd0xpZjdEeFd5a0Z0YitnNTd3Vjg3RnhiczRxNytjVCt1ck1xc2JGbzhNbGE4WWNOUWVXdmtXUk1DU2NHMkpIegp3Zzg9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
        server: https://192.168.0.103:6443
      name: my-cluster
    contexts:
    - context:
        cluster: my-cluster
        namespace: hello
        user: world-hello-admin
      name: my-context
    current-context: my-context
    kind: Config
    preferences: {}
    users:
    - name: world-hello-admin
      user:
        token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJoZWxsbyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJ3b3JsZC10b2tlbi1odDVjbCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJ3b3JsZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjJjMDY4MDIxLWE1ZWYtNGYyMC05M2UyLTc4Y2VkZGIwZDBmMyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpoZWxsbzp3b3JsZCJ9.pb8OAr4-gBB1FlBDaBaySHCEqJtgzp28sMo3XBXk-CWUvcNpmikOi-K-mVlcb8Q13A7EOc1LNKeWKmMnLjq7cRDrMdOGTJ6-ml309VjndQ86INeuOiKvyoL0o0JQp6wdgIRyLLHXjV9X1AJ19CxdtJrpbowPJetooBjZCG1e9G62Pm8JULjh9jomgf3BriHptTcvF-y4tRt6-JHMeQMeDzcnLtOgdI4FweacYanzmF41ZehfDiUuR5XuO9c1R2XkuHUUeEnY7a_k0oyu5ewZQkQKILLoJbPghAI5cuEZTWAfc-3jmIu_Rd2EFEdgrm3pFQuWE4ffjmZYum_9V9iXwg
    [root@test ~]# 


      
     


