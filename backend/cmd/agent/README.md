# Agent process

The agent process runs on a managed Linux server. It is the only component intended to perform privileged local operations, and it must expose typed operations rather than arbitrary shell execution.

Agent operations must validate input, be idempotent where possible, emit structured results, and report observed state back to the control plane.
