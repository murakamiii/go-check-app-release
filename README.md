# go-check-app-release
![Go](https://github.com/murakamiii/go-check-app-release/workflows/Go/badge.svg)

Checking iOS/Android app release on store & post slack messages

```
make run ARG="-slack T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX \
-ios IOSIDXXXXXXX \
-android com.id.android.app \
-cache false\
-register '{{.OS}} {{.Version}} app registered:tada:' \
-update '{{.OS}} {{.Version}} app released:tada:'"
```
