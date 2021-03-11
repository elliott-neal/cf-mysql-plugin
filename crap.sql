FAILED
Incorrect Usage: Unable to parse local forwarding argument: "62130::"

NAME:
   ssh - SSH to an application container instance

USAGE:
   cf ssh APP_NAME [-i app-instance-index] [-c command] [-L [bind_address:]port:host:hostport] [--skip-host-validation] [--skip-remote-execution] [--request-pseudo-tty] [--force-pseudo-tty] [--disable-pseudo-tty]

OPTIONS:
   -L                               Local port forward specification. This flag can be defined more than once.
   --app-instance-index, -i         Application instance index
   --command, -c                    Command to run. This flag can be defined more than once.
   --disable-pseudo-tty, -T         Disable pseudo-tty allocation
   --force-pseudo-tty, -tt          Force pseudo-tty allocation
   --request-pseudo-tty, -t         Request pseudo-tty allocation
   --skip-host-validation, -k       Skip host key validation
   --skip-remote-execution, -N      Do not execute a remote command
