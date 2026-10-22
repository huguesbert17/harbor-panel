# API process

The API process is the control-plane HTTP entry point. It should authenticate requests, authorize module operations, return resource views, and enqueue mutations.

It must not execute privileged server commands directly. Server mutations go through jobs and the managed-server agent.
