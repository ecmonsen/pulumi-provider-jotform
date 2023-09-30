#!/bin/bash
set -euo pipefail
BASE_URL="https://api.jotform.com/docs/properties/index.php"
curl "$BASE_URL" -o index.php.tmp
for OPTION_VALUE in $(grep "option value" index.php.tmp | perl -pe 's# *</select>##;s#<option value=([^ >]+) *>([^<]+)</option>[^<]*#$1\n#g;s/^\s*//;s#<option disabled.*</option>\\n##')
do
  echo "$OPTION_VALUE"
  curl -s "$BASE_URL?field=$OPTION_VALUE" -o "${OPTION_VALUE}.html"
done

