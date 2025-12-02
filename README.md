# Kubernetes Lifecycle Hooks Demo

This project demonstrates Kubernetes container lifecycle hooks (postStart and preStop) using a Go application deployed to minikube.

## What this does

The application implements both lifecycle hooks:

- **postStart**: Logs container start time to `/tmp/poststart.log`
- **preStop**: Continuously logs container stop time to `/tmp/prestop.log` during graceful shutdown

## Test Results

```shell
~/k8s-lifecycle-hooks main
❯ docker build -t test_svc .
[+] Building 0.9s (13/13) FINISHED                                                                                                                                           docker:default

~/k8s-lifecycle-hooks main
❯ minikube image load test_svc:latest

~/k8s-lifecycle-hooks main                                                                                                                                                               4s
❯ minikube image ls | grep test_svc
docker.io/library/test_svc:latest

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ k apply -f deployment.yaml
deployment.apps/test-svc created

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ kpods
NAME                        READY   STATUS    RESTARTS   AGE
test-svc-7bcd5c5697-8lb5p   1/1     Running   0          3s

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ klogs test-svc-7bcd5c5697-8lb5p
Test Service started
PostStart hook output: Container started at Tue Dec  2 15:00:32 UTC 2025

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ k exec -it test-svc-7bcd5c5697-8lb5p -- sh
/ # cat /tmp/poststart.log
Container started at Tue Dec  2 15:00:32 UTC 2025
/ # exit

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ k rollout restart deployment test-svc
deployment.apps/test-svc restarted

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ kpods
NAME                       READY   STATUS        RESTARTS   AGE
test-svc-d8bf66fc4-rrmcm   1/1     Terminating   0          80s

~/k8s-lifecycle-hooks main ?1                                                                                                                                                    󱃾 minikube
❯ klogs test-svc-7bcd5c5697-8lb5p
Test Service started
PostStart hook output: Container started at Tue Dec  2 15:02:42 UTC 2025

Received termination signal, checking preStop hook output...

~/k8s-lifecycle-hooks main !1 ?1                                                                                                                                                 󱃾 minikube
❯ k exec -it test-svc-7bcd5c5697-8lb5p -- sh
/ # tail -f /tmp/p
poststart.log  prestop.log
/ # tail -f /tmp/p
poststart.log  prestop.log
/ # tail -f /tmp/prestop.log
29 Container stopping at Tue Dec  2 15:08:23 UTC 2025
30 Container stopping at Tue Dec  2 15:08:24 UTC 2025
31 Container stopping at Tue Dec  2 15:08:25 UTC 2025
32 Container stopping at Tue Dec  2 15:08:26 UTC 2025
33 Container stopping at Tue Dec  2 15:08:27 UTC 2025
...
92 Container stopping at Tue Dec  2 15:26:26 UTC 2025
93 Container stopping at Tue Dec  2 15:26:27 UTC 2025
94 Container stopping at Tue Dec  2 15:26:28 UTC 2025
95 Container stopping at Tue Dec  2 15:26:29 UTC 2025
96 Container stopping at Tue Dec  2 15:26:30 UTC 2025
97 Container stopping at Tue Dec  2 15:26:31 UTC 2025
98 Container stopping at Tue Dec  2 15:26:32 UTC 2025
99 Container stopping at Tue Dec  2 15:26:33 UTC 2025
100 Container stopping at Tue Dec  2 15:26:34 UTC 2025
command terminated with exit code 137
```

## Summary

1. Built and deployed a Go service with lifecycle hooks to minikube
2. **postStart hook**: Successfully logged container startup time  
3. **preStop hook**: Continuously logged shutdown progress for graceful termination
4. During rollout restart, the preStop hook ran for about 18 minutes before the container was forcefully terminated with exit code 137 (SIGKILL)
