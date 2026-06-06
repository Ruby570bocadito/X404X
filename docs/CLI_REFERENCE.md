# X404X — CLI Reference

## Main Commands

```
x404x [command] [subcommand] [flags]

Commands:
  campaign      Manage red team campaigns
  recon         Reconnaissance operations
  agent         Agent management
  exploit       Exploitation and privilege escalation
  ai            AI assistant (Specter + Apex)
  lateral       Lateral movement
  persistence   Persistence mechanisms
  payload       Payload builder
  listeners     Listener management
  dashboard     Web dashboard control
  db            Database management
  console       msfconsole-style interactive shell
  tui           Bubble Tea terminal UI
  lab           Lab environment control
```

## Campaign

```
x404x campaign start     -t 10.0.0.0/24 -g domain_admin -p stealth
x404x campaign status    [--json]
x404x campaign pause     <campaign_id>
x404x campaign resume    <campaign_id>
x404x campaign report    <campaign_id> [--format json|markdown|pdf]
x404x campaign list      [--status active|completed]
x404x campaign delete    <campaign_id>
```

## Recon

```
x404x recon scan         <target> [--ports 1-1000] [--stealth]
x404x recon osint        <domain> [--github] [--shodan]
x404x recon dns          <domain> [--bruteforce]
x404x recon vuln         <target> [--service all]
```

## Agent

```
x404x agent list         [--campaign <id>] [--status online|dead]
x404x agent interact     <agent_id>
x404x agent tasks        <agent_id> [--list] [--add <command>]
x404x agent kill         <agent_id> [--reason "cleanup"]
x404x agent generate     --os linux --arch amd64 [--stealth]
```

## Exploit

```
x404x exploit scan                        # Scan for local privesc vectors
x404x exploit run        [--vector suid] [--risk safe|medium|high]
x404x exploit cve        <CVE-ID> [--target <ip>]
x404x exploit bruteforce <service> <target>
```

## AI

```
x404x ai chat                           # Interactive Specter chat
x404x ai analyze          <target_data> # Analyze context
x404x ai suggest          [--campaign <id>]  # Get suggestions
x404x ai mode             [auto|manual]
x404x ai model            [list|set <model>]
```

## Lateral Movement

```
x404x lateral scan        [--subnet <cidr>]
x404x lateral propagate   [--method smb|ssh|wmi] [--target <ip>]
x404x lateral relay       [--add <ip:port>] [--chain]
```

## Persistence

```
x404x persistence install [--method cron|ssh|systemd|kernel]
x404x persistence list
x404x persistence remove  <id>
x404x persistence kernel  [load|unload|status]
```

## Dashboard

```
x404x dashboard start     [--port 3000] [--dev]
x404x dashboard stop
x404x dashboard status
```

## Database

```
x404x db migrate          [--up|--down]
x404x db status
x404x db backup           [--output <path>]
x404x db restore          <path>
```

## Lab

```
x404x lab up              [--scenario <name>]
x404x lab down
x404x lab status
x404x lab scenario        [list|load <name>]

Available scenarios:
  - ctf_basic       Basic CTF with 5 targets
  - ad_environment  Active Directory simulation
  - webapp_pentest  Web application testing
  - full_chain      Complete kill chain exercise
```
