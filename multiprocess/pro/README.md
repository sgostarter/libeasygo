
## node

try close child processes after the parent process exit

or, you should call ClearOldProcesses on the parent process start

```bash
	ex, err := os.Executable()
	if err != nil {
		logger.WithFields(l.ErrorField(err)).Error("get executable path failed")

		return
	}

	// ex = strings.ReplaceAll(ex, "_service_", "_")

	item.ProxyItem.Listen, pro.NewWrapper(ex, []string{"-child_listen", item.ProxyItem.Listen}, "",
		item.ProxyItem.Listen, fileLogger)
```
