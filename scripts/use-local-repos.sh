#!/bin/bash
set -e

echo "********************************************************"
echo "*** Update go.mod to point to local repos: client,   ***"
echo "*** gateway, gateway admin, provider admin.          ***"
echo "********************************************************"

REPLACE_TEXT_C="replace github.com/ConsenSys/fc-retrieval-client => ../fc-retrieval-client"
REPLACE_TEXT_CO="replace github.com/ConsenSys/fc-retrieval-commonn => ../fc-retrieval-common"
REPLACE_TEXT_G="replace github.com/ConsenSys/fc-retrieval-gateway => ../fc-retrieval-gateway"
REPLACE_TEXT_GA="replace github.com/ConsenSys/fc-retrieval-gateway-admin => ../fc-retrieval-gateway-admin"
REPLACE_TEXT_PA="replace github.com/ConsenSys/fc-retrieval-provider-admin => ../fc-retrieval-provider-admin"
REPLACE_TEXT_R="replace github.com/ConsenSys/fc-retrieval-register => ../fc-retrieval-register"

sed '/replace .*/d' ../go.mod > ../go.mod.temp
rm ../go.mod
mv ../go.mod.temp ../go.mod

printf "%s\n" "$REPLACE_TEXT_C" >> ../go.mod
printf "%s\n" "$REPLACE_TEXT_CO" >> ../go.mod
printf "%s\n" "$REPLACE_TEXT_G" >> ../go.mod
printf "%s\n" "$REPLACE_TEXT_GA" >> ../go.mod
printf "%s\n" "$REPLACE_TEXT_PA" >> ../go.mod
printf "%s\n" "$REPLACE_TEXT_R" >> ../go.mod

echo go.mod file is now:
cat ../go.mod
