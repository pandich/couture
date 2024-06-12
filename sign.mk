name = couture
sign = "Developer ID Application: Stephen Pandich (D9W3Q7D7N3)"
keyId = "WX54UD585Y"
issuerId = "76e8081d-f643-42b9-a772-11481959ad52"

build = dist/
cmd_amd64 = $(build)$(name)_amd64
cmd_arm64 = $(build)$(name)_arm64
cmd_universal = $(build)$(name)
dmg = $(build)$(name).dmg

.PHONY: signed
signed: dist/couture.stapled.txt

dist/amd64:
	GOARCH=amd64 CGO_ENABLED=1 go build -o $(cmd_amd64) .

dist/arm64:
	GOARCH=arm64 CGO_ENABLED=1 go build -o $(cmd_arm64) .
	lipo -create -output $(cmd_universal) $(cmd_amd64) $(cmd_arm64)
	codesign --deep --force --verbose --sign $(sign) --entitlements .sign/entitlements.mac.plist --options=runtime $(cmd_universal)
	codesign --verify --deep --strict --verbose $(cmd_universal) # Verify the app is signed
	hdiutil create -volname Couture -srcfolder $(cmd_universal) -ov -format UDZO $(dmg)
	xcrun notarytool submit $(dmg) --key .sign/AuthKey_$(keyId).p8 --key-id $(keyId) --issuer $(issuerId) --wait
	xcrun stapler staple $(dmg)
	xcrun stapler validate $(dmg)
	touch dist/couture.stapled.txt
	spctl -a -vvv -t install $(cmd_universal)

dist/universal: dist/amd64 dist/arm64

dist/couture.dmg: dist/universal

dist/couture.stapled.txt: dist/couture.dmg

.PHONY: clean all
clean:
	rm -f $(cmd_amd64) $(cmd_arm64) $(cmd_universal) $(dmg) dist/couture.stapled.txt
	mkdir -p dist

all: clean signed
