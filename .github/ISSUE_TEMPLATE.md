# Please attach the following output when submitting a bug:
- [ ] Attach the value of `cat /proc/sys/kernel/perf_event_paranoid` it should be `0`
- [ ] Attached the value of `uname -a`
- [ ] Attached the value of `mount | grep debugfs`
- [ ] Attached the value of `zcat /proc/config.gz | grep -iE '(perf|debugfs)'` or `cat /boot/config-$(uname -r) | grep -iE '(perf|debugfs)`
- [ ] Attached the build and the version of go runtime `perf_exporter -h`

### Bug reports:
Please replace this line with a brief summary of your issue.
