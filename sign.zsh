#!/usr/bin/env zsh

set -x

go build -o build/couture .
codesign --deep --force --verbose --sign "Developer ID Application: Stephen Pandich (D9W3Q7D7N3)" --entitlements .sign/entitlements.mac.plist --options=runtime  build/couture
hdiutil create -volname Couture -srcfolder build/couture -ov -format UDZO build/couture.dmg
xcrun notarytool submit build/couture.dmg --key .sign/AuthKey_WX54UD585Y.p8 --key-id WX54UD585Y --issuer 76e8081d-f643-42b9-a772-11481959ad52 --wait
xcrun stapler staple build/couture.dmg
xcrun stapler validate build/couture.dmg
