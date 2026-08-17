# cloud2butane

A small prototype transpiler in Go that converts a subset of cloud-init
`cloud-config` YAML into [Butane](https://coreos.github.io/butane/) YAML.

This was built to explore the translation problems involved in the CNCF
Flatcar LFX Mentorship project
("Cloud-Init to Butane YAML config transpiler"), in particular the cases
where cloud-init and Butane look like they map field-to-field but don't,
because cloud-init is imperative and Butane/Ignition is purely declarative.

## What it does

Given a cloud-config file, it produces an equivalent Butane config by:

- Translating `users` (name, groups, shell) into Butane `passwd.users`.
- Translating `write_files` entries into Butane `storage.files`, including
  permissions (parsed as octal and re-emitted in Butane's expected octal
  YAML form).
- Detecting `write_files` entries that write to a systemd unit path
  (`/etc/systemd/system/*.service`, `.timer`, `.socket`, `.mount`,
  `.target`, `.swap`) and translating them into Butane `systemd.units`
  entries instead of plain files.
- Detecting a paired `runcmd` entry (`systemctl enable <unit>` /
  `systemctl disable <unit>`) and using it to set `enabled: true`/`false`
  on the corresponding systemd unit — since Butane/Ignition has no
  equivalent to `runcmd` and can't run arbitrary commands at
  provisioning time. If no matching `runcmd` is found, the unit is
  written but left disabled, matching what the original cloud-config
  would have actually done.
- Detecting `write_files` entries that target a systemd drop-in path
  (`<unit>.d/<name>.conf`) and translating them into a `dropins` entry
  under the corresponding Butane systemd unit, creating a stub unit if
  the base unit wasn't otherwise defined.

## Usage

```sh
go run ./cmd <cloud-config-file> [--debug]
```

`--debug` prints the parsed cloud-config as JSON before translation, useful
for checking how a file was interpreted.

Example:

```sh
go run ./cmd ex_cloud_systemd.yaml
```

See `ex_cloud.yaml`, `ex_cloud_systemd.yaml`, and `ex_cloud_dropin.yaml` for
sample inputs covering plain files, systemd unit fusion, and drop-ins.

## Known limitations / not yet implemented

This is a prototype, scoped to demonstrate the trickiest translation cases
rather than full coverage. Not yet handled:

- **Groups** — cloud-init's top-level `groups` module (creating new groups)
  and Butane's `passwd.groups` aren't parsed. User-to-group *membership* is
  supported; group *creation* isn't yet.
- **Certificates** — cloud-init's `ca_certs` module isn't translated.
- **`append`** — cloud-config's `write_files[].append: true` is parsed but
  not yet reflected in the Butane output; files are currently always
  emitted as a full overwrite via `contents.inline`. Butane's
  `storage.files` supports a separate `append` list for this case.
- **`systemctl enable --now <unit>`** and other multi-flag `systemctl`
  invocations aren't recognized by the enable/disable detection, which
  currently expects a plain `systemctl enable|disable <unit>`.
- Only `runcmd` lines matching `systemctl enable|disable <unit>` are
  understood. Ignition itself has no way to run arbitrary commands at
  provisioning time, so a generic `runcmd` can't be translated directly.
  It might be reasonable to investigate approximating `runcmd` by having Butane declare a
  synthesized systemd oneshot unit (`Type=oneshot`, `ExecStart=<command>`,
  `ConditionFirstBoot=true`) that runs the command at first boot instead —
  this shifts execution from Ignition's provisioning stage to systemd
  itself, after Ignition has finished. That's a reasonable direction for a
  future iteration, but raises design questions (how to name synthesized
  units, whether to merge multiple unmatched `runcmd` lines into one unit
  or one each) that are out of scope for this prototype.

## Tests

```sh
go test ./...
```

Covers systemd unit path detection and the enable/disable detection logic
(including edge cases like repeated/conflicting `systemctl` calls and
extra whitespace).
